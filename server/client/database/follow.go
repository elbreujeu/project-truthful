package database

import (
	"database/sql"
	"log"
)

func CheckFollowExists(follower int, followed int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND followed = ?", follower, followed).Scan(&count)
	if err != nil {
		log.Printf("Error checking if follow exists for follower %d and followed %d, %v\n", follower, followed, err)
		return false, err
	}
	return count > 0, nil
}

func AddFollow(followerId int, followedId int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO follow (follower, followed) VALUES (?, ?)", followerId, followedId)
	if err != nil {
		log.Printf("Error inserting follow for follower %d and followed %d, %v\n", followerId, followedId, err)
		return err
	}
	return nil
}

func RemoveFollow(followerId int, followedId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM follow WHERE follower = ? AND followed = ?", followerId, followedId)
	if err != nil {
		log.Printf("Error deleting follow for follower %d and followed %d, %v\n", followerId, followedId, err)
		return err
	}
	return nil
}
