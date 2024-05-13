package routes

import (
	"log"
	"net/http"
	"os"
	"project_truthful/client/database"
	"project_truthful/client/token"
	"time"

	"github.com/gin-gonic/gin"
)

func SetMiddleware(r *gin.Engine) {
	r.Use(setJSONResponse)
	r.Use(checkAndUpdateRateLimit)
}

func setJSONResponse(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}

func parseAndVerifyAccessToken(c *gin.Context) (int, int, error) {
	accessToken, code, err := token.ParseAccessToken(c)
	if err != nil {
		log.Printf("Error while parsing token: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while parsing token", "error": err.Error()})
		return 0, code, err
	}

	requesterId, code, err := token.VerifyJWT(accessToken)
	if err != nil {
		log.Printf("Error while checking token: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while checking token", "error": err.Error()})
		return 0, code, err
	}
	return requesterId, http.StatusOK, nil
}

func moderationLogging(moderatorId int, action string, targetId int) error {
	err := database.LogModerationAction(moderatorId, action, targetId, database.DB)

	if err != nil {
		log.Printf("Error logging moderation action %s by moderator %d, %v\n", action, moderatorId, err)
	}
	return err
}

// true : rate limit exceeded
// false : rate limit not exceeded
func checkAndUpdateRateLimit(c *gin.Context) {
	// check if rate limit is enforced in env
	rateLimitEnforced := os.Getenv("RATE_LIMIT_ENFORCED")
	if rateLimitEnforced != "true" {
		c.Next()
		return
	}

	// get latest rate limit
	rateLimit, err := database.GetRateLimit(c.ClientIP(), database.DB)

	if err != nil {
		log.Printf("Error getting rate limit for ip %s, %v\n", c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting rate limit", "error": err.Error()})
		c.Abort()
		return
	}

	// if last request was more than 1 hour ago, reset request count
	if rateLimit.LastRequestTime.Add(1 * time.Hour).Before(time.Now()) {
		err = database.ResetRateLimit(c.ClientIP(), database.DB)
		if err != nil {
			log.Printf("Error resetting rate limit for ip %s, %v\n", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error resetting rate limit", "error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
		return
	}

	err = database.IncrementRateLimit(c.ClientIP(), database.DB)
	if err != nil {
		log.Printf("Error incrementing rate limit for ip %s, %v\n", c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error incrementing rate limit", "error": err.Error()})
		c.Abort()
		return
	}

	if rateLimit.RequestCount > 100 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
		c.Abort()
		return
	} else {
		c.Next()
		return
	}
}
