package client

import (
	"errors"
	"net/http"
	"net/mail"
	"project_truthful/client/database"
)

func checkUserInfos(displayName string, email string) error {
	if len(displayName) == 0 {
		return errors.New("display name is empty")
	}
	if len(displayName) > 30 {
		return errors.New("display name is too long")
	}
	if len(email) == 0 {
		return errors.New("email address is empty")
	}
	if len(email) > 319 {
		return errors.New("email address is too long")
	}
	_, err := mail.ParseAddress(email) // todo : check if email is already in db and not the same as the current one
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserInformations(requesterId int, displayName string, email string) (int, error) {
	exists, err := database.CheckUserIdExists(requesterId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exists {
		return http.StatusNotFound, errors.New("user not found")
	}

	err = checkUserInfos(displayName, email)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = database.UpdateUserInformations(requesterId, displayName, email, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
