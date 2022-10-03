package token

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtPublicKey *rsa.PublicKey
var jwtPrivateKey *rsa.PrivateKey

func Init() error {
	pubKey, err := os.ReadFile("./cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}
	jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		log.Printf("Unable to parse ECDSA public key: %v", err)
		return err
	}

	prvKey, err := os.ReadFile("./cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(prvKey))
	if err != nil {
		log.Printf("Unable to parse ECDSA private key: %v", err)
		return err
	}
	return nil
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id":     userID,
		"created_at":  time.Now().Unix(),
		"expiry_date": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtPrivateKey)
}
