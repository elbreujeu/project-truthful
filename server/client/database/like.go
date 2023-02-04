package database

import (
	"database/sql"
	"log"
)

func CheckLikeExists(userId int, postId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer_like WHERE user_id = ? AND answer_id = ?", userId, postId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if like for user %d and post %d exists, %v\n", userId, postId, err)
		return false, err
	}
	return count > 0, nil
}

func RemoveLike(userId int, postId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM answer_like WHERE user_id = ? AND answer_id = ?", userId, postId)
	if err != nil {
		log.Printf("Error deleting like for user %d and post %d, %v\n", userId, postId, err)
		return err
	}
	return nil
}

func AddLike(userId int, postId int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO answer_like (user_id, answer_id) VALUES (?, ?)", userId, postId)
	if err != nil {
		log.Printf("Error inserting like for user %d and post %d, %v\n", userId, postId, err)
		return err
	}
	return nil
}
