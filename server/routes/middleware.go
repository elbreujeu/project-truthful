package routes

import (
	"log"
	"net/http"
	"project_truthful/client/token"

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
