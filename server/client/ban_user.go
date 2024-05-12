package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func BanUser(userId int, requesterId int, duration int, reason string) (int64, int, error) {
	// Checks if the author is an admin
	isModerator, err := database.CheckModeratorStatus(requesterId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	isAdmin, err := database.CheckAdminStatus(requesterId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !isModerator && !isAdmin {
		return 0, http.StatusForbidden, errors.New("user is not a moderator or admin")
	}

	// checks if the user exists
	exists, err := database.CheckUserIdExists(userId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !exists {
		return 0, http.StatusNotFound, errors.New("user not found")
	}

	// checks if user is self
	if userId == requesterId {
		return 0, http.StatusForbidden, errors.New("cannot ban self")
	}

	//checks if user is admin
	isAdmin, err = database.CheckAdminStatus(userId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if isAdmin {
		return 0, http.StatusForbidden, errors.New("cannot ban an admin")
	}

	// Bans the user
	banId, err := database.BanUser(userId, requesterId, duration, reason, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return banId, http.StatusOK, nil
}

func PardonUser(banId int, requesterId int) (int64, int, error) {
	// Checks if the author is an admin
	isModerator, err := database.CheckModeratorStatus(requesterId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	isAdmin, err := database.CheckAdminStatus(requesterId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !isModerator && !isAdmin {
		return 0, http.StatusForbidden, errors.New("user is not a moderator or admin")
	}

	// Checks if the ban exists
	exists, err := database.CheckBanExistsByBanId(banId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if !exists {
		return 0, http.StatusNotFound, errors.New("ban not found")
	}

	// Checks if the ban is already pardoned
	pardoned, err := database.CheckPardonExists(banId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if pardoned {
		return 0, http.StatusForbidden, errors.New("ban is already pardoned")
	}

	// Pardons the user
	pardonId, err := database.PardonUser(banId, requesterId, database.DB)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return pardonId, http.StatusOK, nil
}
