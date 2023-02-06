package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
)

func GetQuestions(userId int, start int, end int, db *sql.DB) ([]models.Question, error) {
	//selects all questions in database where receiver_id = userId
	rows, err := db.Query("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question WHERE receiver_id = ? ORDER BY created_at DESC LIMIT ?, ?", userId, start, end)
	if err != nil {
		log.Printf("Error getting questions for user %d, %v\n", userId, err)
		return nil, err
	}
	defer rows.Close()
	// prints all the rows with their id and text
	var questions []models.Question
	for rows.Next() {
		// TODO : put this in a function
		var curQuestion models.Question
		var authorId sql.NullInt64
		err := rows.Scan(&curQuestion.Id, &curQuestion.Text, &authorId, &curQuestion.IsAuthorAnonymous, &curQuestion.ReceiverId, &curQuestion.CreatedAt)
		if err != nil {
			log.Printf("Error scanning question for user %d, %v\n", userId, err)
			return nil, err
		}
		hasBeenAnswered, err := HasQuestionBeenAnswered(curQuestion.Id, db)
		if err != nil {
			log.Printf("Error checking if question %d has been answered, %v\n", curQuestion.Id, err)
			return nil, err
		}
		if hasBeenAnswered {
			continue
		}
		if curQuestion.IsAuthorAnonymous || !authorId.Valid {
			curQuestion.Author = models.UserPreview{}
		} else {
			curQuestion.Author.Id = authorId.Int64
			authorUsername, authorDisplayName, err := GetUsernameAndDisplayName(int(curQuestion.Author.Id), db)
			if err != nil {
				log.Printf("Error getting author username and display name for question %d, %v\n", curQuestion.Id, err)
				return nil, err
			}
			curQuestion.Author.Username = authorUsername
			curQuestion.Author.DisplayName = authorDisplayName
		}
		questions = append(questions, curQuestion)
	}
	return questions, nil
}

func GetQuestionReceiverId(questionId int, db *sql.DB) (int, error) {
	var userId int
	err := db.QueryRow("SELECT receiver_id FROM question WHERE id = ?", questionId).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting question author for question %d, %v\n", questionId, err)
		return 0, err
	} else if err == sql.ErrNoRows {
		return 0, err
	}
	return userId, nil
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

func GetQuestionById(questionId int, db *sql.DB) (models.Question, error) {
	var question models.Question
	var authorId sql.NullInt64
	err := db.QueryRow("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question WHERE id = ?", questionId).Scan(&question.Id, &question.Text, &authorId, &question.IsAuthorAnonymous, &question.ReceiverId, &question.CreatedAt)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting question %d, %v\n", questionId, err)
		return models.Question{}, err
	} else if err == sql.ErrNoRows {
		return models.Question{}, err
	}
	if question.IsAuthorAnonymous || !authorId.Valid {
		question.Author = models.UserPreview{}
	} else {
		question.Author.Id = authorId.Int64
		authorUsername, authorDisplayName, err := GetUsernameAndDisplayName(int(question.Author.Id), db)
		if err != nil {
			log.Printf("Error getting author username and display name for question %d, %v\n", question.Id, err)
			return models.Question{}, err
		}
		question.Author.Username = authorUsername
		question.Author.DisplayName = authorDisplayName
	}
	return question, nil
}
