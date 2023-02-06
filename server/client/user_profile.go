package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func GetUserProfile(username string, count int, start int) (models.UserProfileInfos, int, error) {
	id, err := database.GetUserId(username, database.DB)
	if err != nil && err != sql.ErrNoRows {
		return models.UserProfileInfos{}, http.StatusInternalServerError, err
	}
	if id == 0 || err == sql.ErrNoRows {
		return models.UserProfileInfos{}, http.StatusNotFound, errors.New("user not found")
	}
	infos, err := database.GetUserProfileInfos(id, count, start, database.DB)
	if err != nil {
		return models.UserProfileInfos{}, http.StatusInternalServerError, err
	}
	return infos, http.StatusOK, nil
}
