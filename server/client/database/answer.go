package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
)

func CheckAnswerIdExists(answerId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer WHERE id = ? AND has_been_deleted = 0", answerId).Scan(&count)
	if err != nil {
		log.Printf("Error checking if post %d exists, %v\n", answerId, err)
		return false, err
	}
	return count > 0, nil
}

func GetAnswerAuthorId(answerId int, db *sql.DB) (int, error) {
	var authorId int
	err := db.QueryRow("SELECT user_id FROM answer WHERE id = ? AND has_been_deleted = 0", answerId).Scan(&authorId)
	if err != nil {
		log.Printf("Error getting author id for answer %d, %v\n", answerId, err)
		return 0, err
	}
	return authorId, nil
}

func HasQuestionBeenAnswered(questionId int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM answer WHERE question_id = ? AND has_been_deleted = 0", questionId).Scan(&count)
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

func MarkAnswerAsDeleted(answerId int, db *sql.DB) error {
	_, err := db.Exec("UPDATE answer SET has_been_deleted = 1, deleted_at = NOW() WHERE id = ?", answerId)
	if err != nil {
		log.Printf("Error marking answer %d as deleted, %v\n", answerId, err)
		return err
	}
	return nil
}

func getAnswers(id int, count int, start int, db *sql.DB) ([]models.Answer, error) {
	rows, err := db.Query("SELECT id, question_id, text, created_at FROM answer WHERE user_id = ? AND has_been_deleted = 0 ORDER BY created_at DESC LIMIT ?, ?", id, start, count+start)
	if err != nil {
		log.Printf("Error getting answers for id %d, %v\n", id, err)
		return nil, err
	}
	defer rows.Close()
	var answers []models.Answer
	for rows.Next() {
		//TODO: Put this in a function
		var answer models.Answer
		var questionId int
		err := rows.Scan(&answer.Id, &questionId, &answer.AnswerText, &answer.CreatedAt)
		if err != nil {
			log.Printf("Error scanning answer for id %d, %v\n", id, err)
			return nil, err
		}
		question, err := GetQuestionById(questionId, db)
		if err != nil {
			log.Printf("Error getting question for answer %d, %v\n", answer.Id, err)
			question = models.Question{}
		}
		answer.QuestionText = question.Text
		if question.IsAuthorAnonymous {
			answer.Author = models.UserPreview{}
			answer.IsAuthorAnonymous = true
		} else {
			answer.Author.Id = question.Author.Id
			answer.Author.Username = question.Author.Username
			answer.Author.DisplayName = question.Author.DisplayName
		}
		answer.LikeCount, err = GetLikeCountForAnswer(answer.Id, db)
		if err != nil {
			log.Printf("Error getting like count for answer %d, %v\n", answer.Id, err)
			answer.LikeCount = 0
		}
		answers = append(answers, answer)
	}
	return answers, nil
}
