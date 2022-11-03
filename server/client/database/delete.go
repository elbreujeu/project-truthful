package database

import (
	"database/sql"
	"log"
)

func RemoveFollow(followerId int, followedId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM follow WHERE follower = ? AND followed = ?", followerId, followedId)
	if err != nil {
		log.Printf("Error deleting follow for follower %d and followed %d, %v\n", followerId, followedId, err)
		return err
	}
	return nil
}
