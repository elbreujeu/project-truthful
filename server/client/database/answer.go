package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
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

func HasQuestionBeenAnswered(questionId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer WHERE question_id = ?", questionId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if question %d has been answered, %v\n", questionId, err)
		return false, err
	}
	return count > 0, nil
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

func getAnswers(id int, start int, end int, db *sql.DB) ([]models.Answer, error) {
	//TODO
	return nil, nil
}
