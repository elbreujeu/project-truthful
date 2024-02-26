package database

import (
	"database/sql"
	"log"
)

func CheckAdminStatus(userId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE id = ? AND is_admin = 1", userId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user %d is an admin, %v\n", userId, err)
		return false, err
	}
	return count > 0, nil
}

func CheckModeratorStatus(userId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE id = ? AND is_moderator = 1", userId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user %d is a moderator, %v\n", userId, err)
		return false, err
	}
	return count > 0, nil
}

func PromoteUserToAdmin(userId int, db *sql.DB) error {
	_, err := db.Exec("UPDATE user SET is_admin = 1 WHERE id = ?", userId)
	if err != nil {
		log.Printf("Error promoting user %d to admin, %v\n", userId, err)
		return err
	}
	return nil
}

func PromoteUserToModerator(userId int, db *sql.DB) error {
	_, err := db.Exec("UPDATE user SET is_moderator = 1 WHERE id = ?", userId)
	if err != nil {
		log.Printf("Error promoting user %d to moderator, %v\n", userId, err)
		return err
	}
	return nil
}
