package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func GetQuestions(userId int, start int, count int) ([]models.Question, int, error) {
	if count > 30 {
		count = 30
	}
	exists, err := database.CheckUserIdExists(userId, database.DB)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !exists {
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	questions, err := database.GetQuestions(userId, start, start+count, database.DB)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return questions, 0, nil
}
