package client

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"project_truthful/client/database"
	"project_truthful/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func isUsernameValid(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if database.CheckUsernameExists(username) {
		return errors.New("username already exists")
	}
	return nil
}

func isPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	hasNumber := false
	hasLowercase := false
	hasUppercase := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
		} else if char >= 'a' && char <= 'z' {
			hasLowercase = true
		} else if char >= 'A' && char <= 'Z' {
			hasUppercase = true
		}
	}
	if !hasNumber || !hasLowercase || !hasUppercase {
		return errors.New("password must contain at least one number, one lowercase letter, and one uppercase letter")
	}
	return nil
}

func isEmailValid(email string) error {
	if len(email) > 319 {
		return errors.New("email must 319 characters at most")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	if database.CheckEmailExists(email) {
		return errors.New("email already exists")
	}
	return nil
}

func isBirthdateValid(birthdateStr string) error {
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		return err
	}
	//checks if birthdate is in the future
	if birthdate.After(time.Now()) {
		return errors.New("birthdate cannot be in the future")
	}
	//checks if birthdate is more than 13 years ago
	if birthdate.After(time.Now().AddDate(-13, 0, 0)) {
		return errors.New("birthdate must be more than 13 years ago")
	}
	return nil
}

func CreateUser(infos models.CreateUserInfos) (int64, int, error) {
	log.Printf("Creating user %s\n", infos.Username)

	err := isUsernameValid(infos.Username)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}
	err = isPasswordValid(infos.Password)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}
	err = isEmailValid(infos.Email)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}
	err = isBirthdateValid(infos.Birthdate)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	encryptedPassword, err := encryptPassword(infos.Password)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := database.InsertUser(infos.Username, encryptedPassword, infos.Email, infos.Birthdate)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	log.Printf("User %s created with id %d\n", infos.Username, id)
	return id, http.StatusOK, nil
}
