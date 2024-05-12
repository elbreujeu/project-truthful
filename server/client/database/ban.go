package database

import (
	"database/sql"
	"log"
	"time"
)

func BanUser(userId int, requesterId int, duration int, reason string, db *sql.DB) (int64, error) {
	var banId int64
	if duration == 0 {
		result, err := db.Exec("INSERT INTO ban (user_id, author_id, reason) VALUES (?, ?, ?)", userId, requesterId, reason)
		if err != nil {
			log.Printf("Error banning user %d, %v\n", userId, err)
			return 0, err
		}
		banId, err = result.LastInsertId()
		if err != nil {
			log.Printf("Error getting ban ID, %v\n", err)
			return 0, err
		}
	} else {
		expiration := time.Now().Add(time.Duration(duration) * time.Hour)
		result, err := db.Exec("INSERT INTO ban (user_id, author_id, reason, expires_at) VALUES (?, ?, ?, ?)", userId, requesterId, reason, expiration)
		if err != nil {
			log.Printf("Error banning user %d, %v\n", userId, err)
			return 0, err
		}
		banId, err = result.LastInsertId()
		if err != nil {
			log.Printf("Error getting ban ID, %v\n", err)
			return 0, err
		}
	}
	return banId, nil
}

func CheckUserBanStatus(userId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ban WHERE user_id = ? AND expires_at > NOW()", userId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user %d is banned, %v\n", userId, err)
		return false, err
	}
	return count > 0, nil
}

func CheckBanExistsByBanId(banId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ban WHERE id = ? AND expires_at > NOW()", banId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if ban %d exists, %v\n", banId, err)
		return false, err
	}
	return count > 0, nil
}
