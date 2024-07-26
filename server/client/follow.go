package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/models"
)

func FollowUser(followerId int, followeeId int) (int, error) {
	if followeeId == followerId {
		return http.StatusBadRequest, errors.New("user can't follow himself")
	}

	followeeExists, err := database.CheckUserIdExists(followeeId, database.DB)

	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !followeeExists {
		return http.StatusNotFound, errors.New("followee not found")
	}

	followerExists, err := database.CheckUserIdExists(followerId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !followerExists {
		return http.StatusNotFound, errors.New("follower not found")
	}

	followExists, err := database.CheckFollowExists(followerId, followeeId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if followExists {
		return http.StatusBadRequest, errors.New("user already follows this user")
	}

	err = database.AddFollow(followerId, followeeId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func UnfollowUser(followerId int, followeeId int) (int, error) {
	followExists, err := database.CheckFollowExists(followerId, followeeId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !followExists {
		return http.StatusBadRequest, errors.New("user doesn't follow this user")
	}

	err = database.RemoveFollow(followerId, followeeId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetFollowers(username string, count int, start int) ([]models.UserPreview, int, error) {
	userId, err := database.GetUserId(username, database.DB)
	if err != nil && err == sql.ErrNoRows {
		return []models.UserPreview{}, http.StatusNotFound, errors.New("user not found")
	} else if err != nil {
		return []models.UserPreview{}, http.StatusInternalServerError, err
	}

	followers, err := database.GetFollowers(userId, count, start, database.DB)
	if err != nil {
		return []models.UserPreview{}, http.StatusInternalServerError, err
	}

	var followPreview []models.UserPreview
	for _, follower := range followers {
		username, displayName, err := database.GetUsernameAndDisplayName(follower, database.DB)
		if err != nil {
			return []models.UserPreview{}, http.StatusInternalServerError, err
		}
		followPreview = append(followPreview, models.UserPreview{Id: int64(follower), Username: username, DisplayName: displayName})
	}

	return followPreview, http.StatusOK, nil
}
