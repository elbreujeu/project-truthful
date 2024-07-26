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

func GetFollowers(userId int, count int, start int, db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT follower FROM follow WHERE followed = ? ORDER BY id DESC LIMIT ? OFFSET ?", userId, count, start)
	if err != nil {
		log.Printf("Error getting followers for user %d, %v\n", userId, err)
		return []int{}, err
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var follower int
		err := rows.Scan(&follower)
		if err != nil {
			log.Printf("Error scanning follower for user %d, %v\n", userId, err)
			return []int{}, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}
