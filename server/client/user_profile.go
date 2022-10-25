package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func GetUserProfile(username string) (models.UserProfileInfos, int, error) {
	id, err := database.GetUserId(username, database.DB)
	if err != nil {
		return models.UserProfileInfos{}, http.StatusInternalServerError, err
	}
	if id == 0 {
		return models.UserProfileInfos{}, http.StatusNotFound, errors.New("user not found")
	}
	infos, err := database.GetUserProfileInfos(id, database.DB)
	if err != nil {
		return models.UserProfileInfos{}, http.StatusInternalServerError, err
	}
	return infos, http.StatusOK, nil
}
