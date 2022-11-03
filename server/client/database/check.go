package database

import (
	"database/sql"
	"log"
)

func CheckAnswerIdExists(answerId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer WHERE id = ?", answerId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if post %d exists, %v\n", answerId, err)
		return false, err
	}
	return count > 0, nil
}

func CheckLikeExists(userId int, postId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer_like WHERE user_id = ? AND answer_id = ?", userId, postId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if like for user %d and post %d exists, %v\n", userId, postId, err)
		return false, err
	}
	return count > 0, nil
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

func CheckFollowExists(follower int, followed int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower = ? AND followed = ?", follower, followed).Scan(&count)
	if err != nil {
		log.Printf("Error checking if follow exists for follower %d and followed %d, %v\n", follower, followed, err)
		return false, err
	}
	return count > 0, nil
}

func HasQuestionBeenAnswered(questionId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer WHERE question_id = ?", questionId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if question %d has been answered, %v\n", questionId, err)
		return false, err
	}
	return count > 0, nil
}
