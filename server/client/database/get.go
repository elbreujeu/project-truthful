package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
)

func GetUserId(username string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting user id for username %s, %v\n", username, err)
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	}
	return id, nil
}

func GetHashedPassword(id int, db *sql.DB) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&password)
	if err != nil {
		log.Printf("Error getting hashed password for id %d, %v\n", id, err)
		return "", err
	}
	return password, nil
}

func getAnswers(id int, start int, end int, db *sql.DB) ([]models.Answer, error) {
	//TODO
	return nil, nil
}

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

func GetQuestions(userId int, start int, end int, db *sql.DB) ([]models.Question, error) {
	//selects all questions in database where receiver_id = userId
	rows, err := db.Query("SELECT id, text, author_id, is_author_anonymous, receiver_id, creation_date FROM question WHERE receiver_id = ? ORDER BY creation_date DESC LIMIT ?, ?", userId, start, end)
	if err != nil {
		log.Printf("Error getting questions for user %d, %v\n", userId, err)
		return nil, err
	}
	defer rows.Close()
	// prints all the rows with their id and text
	var questions []models.Question
	for rows.Next() {
		var curQuestion models.Question
		var authorId sql.NullInt64
		err := rows.Scan(&curQuestion.Id, &curQuestion.Text, &authorId, &curQuestion.IsAuthorAnonymous, &curQuestion.ReceiverId, &curQuestion.CreatedAt)
		if err != nil {
			log.Printf("Error scanning question for user %d, %v\n", userId, err)
			return nil, err
		}
		if authorId.Valid {
			curQuestion.AuthorId = authorId.Int64
		} else {
			curQuestion.AuthorId = 0
		}
		hasBeenAnswered, err := HasQuestionBeenAnswered(curQuestion.Id, db)
		if err != nil {
			log.Printf("Error checking if question %d has been answered, %v\n", curQuestion.Id, err)
			return nil, err
		}
		if !hasBeenAnswered {
			questions = append(questions, curQuestion)
		}
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
