package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/client/token"
	"project_truthful/models"

	"golang.org/x/crypto/bcrypt"
)

func Login(infos models.LoginInfos) (string, int, error) {
	id, err := database.GetUserId(infos.Username, database.DB)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if id == 0 {
		return "", http.StatusBadRequest, errors.New("username does not exist")
	}
	hashedPassword, err := database.GetHashedPassword(infos.Username, database.DB)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(infos.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", http.StatusBadRequest, errors.New("invalid login credentials. Please try again")
	} else if err != nil {
		return "", http.StatusInternalServerError, err
	}

	accessToken, err := token.GenerateJWT(id)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return accessToken, http.StatusOK, nil
}
