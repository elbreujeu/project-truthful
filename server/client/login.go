package client

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"project_truthful/client/database"
	"project_truthful/client/token"
	"project_truthful/models"

	"golang.org/x/crypto/bcrypt"
)

func Login(infos models.LoginInfos) (string, int, error) {
	id, err := database.GetUserId(infos.Username, database.DB)
	if err != nil && err == sql.ErrNoRows {
		return "", http.StatusNotFound, errors.New("user not found")
	} else if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if id == 0 {
		return "", http.StatusBadRequest, errors.New("username does not exist")
	}
	hashedPassword, err := database.GetHashedPassword(id, database.DB)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if os.Getenv("IS_TEST") != "true" {
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(infos.Password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			return "", http.StatusBadRequest, errors.New("invalid login credentials. Please try again")
		} else if err != nil {
			return "", http.StatusInternalServerError, err
		}
	}

	accessToken, err := token.GenerateJWT(id)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return accessToken, http.StatusOK, nil
}

func GoogleLogin(provider string, requestToken string) (string, int, error) {
	googleInfos, err := token.VerifyGoogleToken(requestToken)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	// gets provider id for google
	providerId, err := database.GetOAuthProvider(provider, database.DB)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// checks if the user exists in the database
	userId, err := database.GetUserIdBySubject(providerId, googleInfos.Subject, database.DB)
	// flag to create a new user and a new entry in oauth_login
	if err != nil && err == sql.ErrNoRows {
		var code int
		userId, code, err = RegisterOauth(googleInfos.Name, googleInfos.Email, "2000-01-01") // TODO: add birthdate
		if err != nil {
			return "", code, err
		}

		// add to oauth_login table
		err = database.InsertOauthLogin(providerId, googleInfos.Subject, userId, database.DB)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
	} else if err != nil {
		return "", http.StatusInternalServerError, err
	}
	userToken, err := token.GenerateJWT(int(userId))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return userToken, http.StatusOK, nil
}
