package database

import (
	"database/sql"
	"log"
)

func AddFollow(followerId int, followedId int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO follow (follower, followed) VALUES (?, ?)", followerId, followedId)
	if err != nil {
		log.Printf("Error inserting follow for follower %d and followed %d, %v\n", followerId, followedId, err)
		return err
	}
	return nil
}

func AddQuestion(question string, authorId int, authorIpAddress string, isAuthorAnonymous bool, receiverId int, db *sql.DB) (int64, error) {
	var result sql.Result
	var err error

	if authorId == 0 {
		result, err = db.Exec("INSERT INTO question (text, author_ip_address, receiver_id) VALUES (?, ?, ?)", question, authorIpAddress, receiverId)
	} else {
		result, err = db.Exec("INSERT INTO question (text, author_id, author_ip_address, is_author_anonymous, receiver_id) VALUES (?, ?, ?, ?, ?)", question, authorId, authorIpAddress, isAuthorAnonymous, receiverId)
	}
	if err != nil {
		log.Printf("Error inserting question for author %d and receiver %d, %v\n", authorId, receiverId, err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted id for question, %v\n", err)
		return 0, err
	}
	return id, nil
}

func AddAnswer(userId int, questionId int, answerText string, answererIpAddress string, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO answer (user_id, question_id, text, answerer_ip_address) VALUES (?, ?, ?, ?)", userId, questionId, answerText, answererIpAddress)
	if err != nil {
		log.Printf("Error inserting answer for user %d and question %d, %v\n", userId, questionId, err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted id for question, %v\n", err)
		return 0, err
	}
	return id, nil
}

func AddLike(userId int, postId int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO answer_like (user_id, answer_id) VALUES (?, ?)", userId, postId)
	if err != nil {
		log.Printf("Error inserting like for user %d and post %d, %v\n", userId, postId, err)
		return err
	}
	return nil
}

func InsertUser(username string, password string, email string, birthdate string, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO user (username, display_name, password, email, birthdate) VALUES (?, ?, ?, ?, ?)", username, username, password, email, birthdate)
	if err != nil {
		log.Printf("Error inserting user %s, %v\n", username, err)
		return 0, err
	}
	return result.LastInsertId()
}
