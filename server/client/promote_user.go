package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func PromoteUser(requesterId int, userId int, promoteType string) (int, error) {
	// Check if the requester is an admin
	isUserAdmin, err := database.CheckAdminStatus(requesterId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !isUserAdmin {
		return http.StatusForbidden, errors.New("user has no permission to promote users")
	}

	// Promote the user
	if promoteType == "admin" {
		// Check if the user is already an admin
		isUserAdmin, err := database.CheckAdminStatus(userId, database.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if isUserAdmin {
			return http.StatusBadRequest, errors.New("user is already an admin")
		}
		// Promote the user to admin
		err = database.PromoteUserToAdmin(userId, database.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	} else if promoteType == "moderator" {
		// Check if the user is already a moderator
		isUserModerator, err := database.CheckModeratorStatus(userId, database.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if isUserModerator {
			return http.StatusBadRequest, errors.New("user is already a moderator")
		}
		// Promote the user to moderator
		err = database.PromoteUserToModerator(userId, database.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	} else {
		return http.StatusBadRequest, errors.New("invalid promote type")
	}
}
