package routes

import (
	"log"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/client/token"
	"time"

	"github.com/gin-gonic/gin"
)

func SetMiddleware(r *gin.Engine) {
	r.Use(setCORS)
	r.Use(setJSONResponse)
}

func setCORS(c *gin.Context) {
	// note: maybe not secured
	// TODO: check this
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", `"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"`)
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Next()
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
func checkAndUpdateRateLimit(userIp string) (bool, error) {
	// get latest rate limit
	rateLimit, err := database.GetRateLimit(userIp, database.DB)

	if err != nil {
		log.Printf("Error getting rate limit for ip %s, %v\n", userIp, err)
		return false, err
	}

	// if last request was more than 1 hour ago, reset request count
	if rateLimit.LastRequestTime.Add(1 * time.Hour).Before(time.Now()) {
		err = database.ResetRateLimit(userIp, database.DB)
		if err != nil {
			log.Printf("Error resetting rate limit for ip %s, %v\n", userIp, err)
			return false, err
		}
		return false, nil
	}

	err = database.IncrementRateLimit(userIp, database.DB)
	if err != nil {
		log.Printf("Error incrementing rate limit for ip %s, %v\n", userIp, err)
		return false, err
	}

	if rateLimit.RequestCount > 100 {
		return true, nil
	} else {
		return false, nil
	}
}
