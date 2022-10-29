package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func checkQuestionInfos(question string) error {
	if len(question) == 0 {
		return errors.New("question is empty")
	}
	if len(question) > 500 {
		return errors.New("question is too long")
	}
	return nil
}

func AskQuestion(question string, authorId int, authorIpAddress string, isAuthorAnonymous bool, receiverId int) (int64, int, error) {
	receiverExists, err := database.CheckUserIdExists(receiverId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !receiverExists {
		return 0, http.StatusNotFound, errors.New("receiver not found")
	}

	err = checkQuestionInfos(question)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	id, err := database.AddQuestion(question, authorId, authorIpAddress, isAuthorAnonymous, receiverId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}
