package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
)

func GetUserProfileInfos(id int, db *sql.DB) (models.UserProfileInfos, error) {
	var username string
	var displayName string
	err := db.QueryRow("SELECT username, display_name FROM user WHERE id = ?", id).Scan(&username, &displayName)
	if err != nil {
		log.Printf("Error getting user profile infos for id %d, %v\n", id, err)
		return models.UserProfileInfos{}, err
	}

	var followerCount int
	err = db.QueryRow("SELECT COUNT(*) FROM follow WHERE followed = ?", id).Scan(&followerCount)
	if err != nil {
		log.Printf("Error getting followers count for id %d, %v\n", id, err)
		return models.UserProfileInfos{}, err
	}

	var followingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ?", id).Scan(&followingCount)
	if err != nil {
		log.Printf("Error getting following count for id %d, %v\n", id, err)
		return models.UserProfileInfos{}, err
	}

	var answerCount int
	err = db.QueryRow("SELECT COUNT(*) FROM answer WHERE user_id = ?", id).Scan(&answerCount)
	if err != nil {
		log.Printf("Error getting answer count for id %d, %v\n", id, err)
		return models.UserProfileInfos{}, err
	}

	answers, err := getAnswers(id, 0, 10, db)
	if err != nil {
		log.Printf("Error getting answers for id %d, %v\n", id, err)
		return models.UserProfileInfos{}, err
	}

	return models.UserProfileInfos{Id: id, Username: username, DisplayName: displayName, FollowerCount: followerCount, FollowingCount: followingCount, AnswerCount: answerCount, Answers: answers}, nil
}