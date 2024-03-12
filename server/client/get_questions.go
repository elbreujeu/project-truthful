package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func GetQuestions(userId int, start int, count int) ([]models.Question, int, error) {
	if count < 0 || count > 30 {
		count = 30
	}
	if start < 0 {
		start = 0
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
	return questions, http.StatusOK, nil
}

func ModerationGetUserQuestions(requesterId int, username string, start int, count int) ([]models.Question, int, error) {
	isModerator, err := database.CheckModeratorStatus(requesterId, database.DB)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	isAdmin, err := database.CheckAdminStatus(requesterId, database.DB)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !isModerator && !isAdmin {
		return nil, http.StatusForbidden, errors.New("user is not a moderator or admin")
	}

	userId, err := database.GetUserId(username, database.DB)
	if err == sql.ErrNoRows {
		return nil, http.StatusNotFound, errors.New("user not found")
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return GetQuestions(userId, start, start+count)
}
