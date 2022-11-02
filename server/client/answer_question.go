package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func checkAnswerInfos(answer string) error {
	if len(answer) == 0 {
		return errors.New("answer is empty")
	}
	if len(answer) > 1000 {
		return errors.New("answer is too long")
	}
	return nil
}

func AnswerQuestion(userId int, questionId int, answerText string, authorIpAddress string) (int64, int, error) {
	err := checkAnswerInfos(answerText)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	userExists, err := database.CheckUserIdExists(userId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !userExists {
		return 0, http.StatusNotFound, errors.New("user not found")
	}

	questionReceiverId, err := database.GetQuestionReceiverId(questionId, database.DB)
	if err != nil && err == sql.ErrNoRows {
		// if the question doesn't exist, we return a 404
		return 0, http.StatusNotFound, errors.New("question not found")
	} else if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	if questionReceiverId != userId {
		return 0, http.StatusForbidden, errors.New("user is not the receiver of the question")
	}

	// we check if the user has already answered the question
	alreadyAnswered, err := database.HasQuestionBeenAnswered(questionId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if alreadyAnswered {
		return 0, http.StatusForbidden, errors.New("user has already answered the question")
	}

	id, err := database.AddAnswer(userId, questionId, answerText, authorIpAddress, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return id, http.StatusCreated, nil
}
