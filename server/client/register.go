package client

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"project_truthful/client/basicfuncs"
	"project_truthful/client/database"
	"project_truthful/models"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	if os.Getenv("IS_TEST") == "true" {
		return password, nil
	}
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
	usernameExists, err := database.CheckUsernameExists(username, database.DB)
	if err != nil {
		return err
	} else if usernameExists {
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
	emailExists, err := database.CheckEmailExists(email, database.DB)
	if err != nil {
		return err
	} else if emailExists {
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

func Register(infos models.RegisterInfos) (int64, int, error) {
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

	id, err := database.InsertUser(infos.Username, encryptedPassword, infos.Email, infos.Birthdate, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	log.Printf("User %s created with id %d\n", infos.Username, id)
	return id, http.StatusCreated, nil
}

func RegisterOauth(username string, email string, birthdate string) (int64, int, error) {
	// TODO : write tests
	log.Printf("Creating user %s\n", username)
	isTest := os.Getenv("IS_TEST") == "true"

	var usernameForDb string = basicfuncs.DeleteNonAlphanumeric(username)
	err := isUsernameValid(usernameForDb)
	if err != nil {
		// creates a new random username
		for i := 0; isUsernameValid(usernameForDb) != nil && i < 1000; i++ {
			usernameForDb = "user-" + strconv.Itoa(rand.Intn(1000000))
			if isTest {
				usernameForDb = "user-1000"
			}
		}
	}

	err = isEmailValid(email)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	err = isBirthdateValid(birthdate)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	// if display name length is more than 20 characters, it is truncated, if it is empty, it is set to the username
	var displayName string
	if len(username) > 20 {
		displayName = username[:20]
	} else if len(username) == 0 {
		displayName = usernameForDb
	} else {
		displayName = username
	}
	id, err := database.InsertUserWithDisplayName(usernameForDb, displayName, "", email, birthdate, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	log.Printf("User %s created with id %d\n", username, id)
	return id, http.StatusCreated, nil
}
