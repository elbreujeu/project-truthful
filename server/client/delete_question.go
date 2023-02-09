package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func MarkQuestionAsDeleted(userId int, questionId int) (int, error) {
	authorId, err := database.GetQuestionReceiverId(questionId, database.DB)
	if err != nil && err == sql.ErrNoRows {
		// if the question doesn't exist, we return a 404
		return http.StatusNotFound, errors.New("question not found")
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if authorId != userId {
		return http.StatusForbidden, errors.New("user is not the receiver of the question")
	}

	// we check if the user has already answered the question. if so, we delete the answer
	// below code is shitty. we should get the answer id from the db
	answerId, err := database.GetAnswerIdByQuestionId(questionId, database.DB)
	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, err
	}
	if err != sql.ErrNoRows {
		err = database.MarkAnswerAsDeleted(answerId, database.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	err = database.MarkQuestionAsDeleted(questionId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}