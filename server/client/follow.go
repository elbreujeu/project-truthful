package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func FollowUser(followerId int, followeeId int) (int, error) {
	if followeeId == followerId {
		return http.StatusBadRequest, errors.New("user can't follow himself")
	}

	followeeExists, err := database.CheckUserExists(followeeId, database.DB)

	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !followeeExists {
		return http.StatusNotFound, errors.New("followee not found")
	}

	followerExists, err := database.CheckUserExists(followerId, database.DB)
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
