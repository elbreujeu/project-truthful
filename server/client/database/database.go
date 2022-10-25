package database

import (
	"database/sql"
	"log"
	"os"
	"project_truthful/models"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	log.Println("Connecting to db...")
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_CONTAINER_NAME") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Printf("SQL database open error, %v\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("SQL database ping error, %v\n", err)
		return nil, err
	}
	log.Println("Connected!")
	return db, nil
}

func InsertUser(username string, password string, email string, birthdate string, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO user (username, display_name, password, email, birthdate) VALUES (?, ?, ?, ?, ?)", username, username, password, email, birthdate)
	if err != nil {
		log.Printf("Error inserting user %s, %v\n", username, err)
		return 0, err
	}
	return result.LastInsertId()
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

func GetUserId(username string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting user id for username %s, %v\n", username, err)
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, nil
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

func CheckEmailExists(email string, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if email %s exists, %v\n", email, err)
		return false, err
	}
	return count > 0, nil
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

func AddQuestion(question string, authorId int, authorIpAddress string, receiverId int, db *sql.DB) (int64, error) {
	var result sql.Result
	var err error

	if authorId == 0 {
		result, err = db.Exec("INSERT INTO question (text, author_ip_address, receiver_id) VALUES (?, ?, ?)", question, authorIpAddress, receiverId)
	} else {
		result, err = db.Exec("INSERT INTO question (text, author_id, author_ip_address, receiver_id) VALUES (?, ?, ?, ?)", question, authorId, authorIpAddress, receiverId)
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
