package database

import (
	"database/sql"
	"log"
)

func PardonUser(banId int, requesterId int, db *sql.DB) (int64, error) {
	var pardonId int64

	result, err := db.Exec("INSERT INTO pardon (ban_id, pardoner_id) VALUES (?, ?)", banId, requesterId)
	if err != nil {
		log.Printf("Error pardoning user %d, %v\n", banId, err)
		return 0, err
	}

	pardonId, err = result.LastInsertId()
	if err != nil {
		log.Printf("Error getting pardon ID, %v\n", err)
		return 0, err
	}

	return pardonId, nil
}

func CheckPardonExists(banId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pardon WHERE ban_id = ?", banId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if ban %d is pardoned, %v\n", banId, err)
		return false, err
	}
	return count > 0, nil
}
