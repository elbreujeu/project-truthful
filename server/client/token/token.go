package token

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"project_truthful/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/idtoken"
)

var jwtPublicKey *rsa.PublicKey
var jwtPrivateKey *rsa.PrivateKey

func Init() error {
	pubKey, err := os.ReadFile("/cert/id_rsa.pub")
	if err != nil {
		log.Printf("Unable to read public key: %v", err)
		return err
	}
	jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		log.Printf("Unable to parse RSA public key: %v", err)
		return err
	}

	prvKey, err := os.ReadFile("/cert/id_rsa")
	if err != nil {
		log.Printf("Unable to read private key: %v", err)
		return err
	}
	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(prvKey))
	if err != nil {
		log.Printf("Unable to parse RSA private key: %v", err)
		return err
	}
	return nil
}

func GenerateJWT(userID int) (string, error) {
	if os.Getenv("IS_TEST") == "true" {
		return "test", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id":     userID,
		"created_at":  time.Now().Unix(),
		"expiry_date": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtPrivateKey)
}

func VerifyJWT(tokenString string) (int, int, error) {
	if os.Getenv("IS_TEST") == "true" {
		return 1, http.StatusOK, nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtPublicKey, nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return 0, http.StatusInternalServerError, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiryDate := int64(claims["expiry_date"].(float64))
		if expiryDate < time.Now().Unix() {
			return 0, http.StatusBadRequest, errors.New("token expired")
		}
		return int(claims["user_id"].(float64)), http.StatusAccepted, nil
	}
	return 0, http.StatusBadRequest, errors.New("invalid token")
}

func RefreshJWT(tokenString string) (string, int, error) {
	if os.Getenv("IS_TEST") == "true" {
		return "test", http.StatusOK, nil
	}
	id, status, err := VerifyJWT(tokenString)
	if err != nil {
		return "", status, err
	}
	newToken, err := GenerateJWT(id)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return newToken, http.StatusOK, nil
}

func ParseAccessToken(c *gin.Context) (string, int, error) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" || len(accessToken) < 7 || accessToken[:7] != "Bearer " {
		return "", http.StatusBadRequest, errors.New("missing fields")
	}
	return accessToken[7:], http.StatusOK, nil
}

func VerifyGoogleToken(requestToken string) (models.GoogleInfos, error) {
	clientId := os.Getenv("REACT_APP_GOOGLE_CLIENT_ID")
	isTest := os.Getenv("IS_TEST") == "true"

	if isTest {
		return models.GoogleInfos{Name: "toto123", Email: "toto123@gmail.com", Subject: "123456"}, nil
	}

	// Verify the token
	payload, err := idtoken.Validate(context.Background(), requestToken, clientId)
	if err != nil {
		log.Printf("Unable to verify token: %v", err)
		return models.GoogleInfos{}, err
	}

	// Extract the claims
	var googleInfos models.GoogleInfos
	claims := payload.Claims

	googleInfos.Issuer, _ = claims["iss"].(string)
	googleInfos.AuthorizedParty, _ = claims["azp"].(string)
	googleInfos.Audience, _ = claims["aud"].(string)
	googleInfos.Subject, _ = claims["sub"].(string)
	googleInfos.Email, _ = claims["email"].(string)
	googleInfos.EmailVerified, _ = claims["email_verified"].(bool)
	googleInfos.NotValidBefore, _ = claims["nbf"].(int64)
	googleInfos.Name, _ = claims["name"].(string)
	googleInfos.Picture, _ = claims["picture"].(string)
	googleInfos.GivenName, _ = claims["given_name"].(string)
	googleInfos.FamilyName, _ = claims["family_name"].(string)
	googleInfos.IssuedAt, _ = claims["iat"].(int64)
	googleInfos.ExpiryDate, _ = claims["exp"].(int64)
	googleInfos.JwtId, _ = claims["jti"].(string)

	// prints the whole claims
	fmt.Printf("Claims: %v\n", claims) // debug
	return googleInfos, nil
}
