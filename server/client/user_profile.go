package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func GetUserProfile(username string, requestingUserId int, count int, start int) (models.UserProfileInfos, int, error) {
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

	// If the user requesting the profile is present, check if they are the same user or if they are following the user
	if requestingUserId != 0 {
		if requestingUserId == id {
			infos.IsRequestingSelf = true
		} else {
			infos.IsRequestingSelf = false
			isFollowedByRequester, err := database.CheckFollowExists(requestingUserId, id, database.DB)
			if err != nil {
				return models.UserProfileInfos{}, http.StatusInternalServerError, err
			}
			infos.IsFollowedByRequester = isFollowedByRequester
		}
	} else {
		infos.IsRequestingSelf = false
		infos.IsFollowedByRequester = false
	}

	return infos, http.StatusOK, nil
}
