package database

import (
	"database/sql"
	"log"
)

func InsertUser(username string, password string, email string, birthdate string, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO user (username, display_name, password, email, birthdate) VALUES (?, ?, ?, ?, ?)", username, username, password, email, birthdate)
	if err != nil {
		log.Printf("Error inserting user %s, %v\n", username, err)
		return 0, err
	}
	return result.LastInsertId()
}

func CheckUsernameExists(username string, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Printf("Error checking if username %s exists, %v\n", username, err)
		return false, err
	}
	return count > 0, nil
}

func CheckUserIdExists(id int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE id = ?", id).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user %d exists, %v\n", id, err)
		return false, err
	}
	return count > 0, nil
}

func CheckEmailExists(email string, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if email %s exists, %v\n", email, err)
		return false, err
	}
	return count > 0, nil
}

func GetUserId(username string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting user id for username %s, %v\n", username, err)
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	}
	return id, nil
}

func GetHashedPassword(id int, db *sql.DB) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&password)
	if err != nil {
		log.Printf("Error getting hashed password for id %d, %v\n", id, err)
		return "", err
	}
	return password, nil
}

func GetUsernameAndDisplayName(id int, db *sql.DB) (string, string, error) {
	var username string
	var displayName string
	err := db.QueryRow("SELECT username, display_name FROM user WHERE id = ?", id).Scan(&username, &displayName)
	if err != nil {
		log.Printf("Error getting username and display name for id %d, %v\n", id, err)
		return "", "", err
	}
	return username, displayName, nil
}

func UpdateUserInformations(id int, displayName string, email string, db *sql.DB) error {
	_, err := db.Exec("UPDATE user SET display_name = ?, email = ? WHERE id = ?", displayName, email, id)
	if err != nil {
		log.Printf("Error updating user informations for id %d, %v\n", id, err)
		return err
	}
	return nil
}
