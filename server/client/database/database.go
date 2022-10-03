package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	log.Println("Connecting to db...")
	var db *sql.DB
	var err error
	if os.Getenv("IS_TEST") != "true" {
		cfg := mysql.Config{
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PASSWORD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DB_CONTAINER_NAME") + ":" + os.Getenv("DB_PORT"),
			DBName:               os.Getenv("DB_NAME"),
			AllowNativePasswords: true,
			ParseTime:            true,
		}
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Printf("SQL database open error, %v\n", err)
			return nil, err
		}
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

func CheckUsernameExists(username string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Printf("Error checking if username %s exists, %v\n", username, err)
		return false
	}
	return count > 0
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

func GetHashedPassword(username string, db *sql.DB) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&password)
	if err != nil {
		log.Printf("Error getting hashed password for username %s, %v\n", username, err)
		return "", err
	}
	return password, nil
}

func CheckEmailExists(email string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if email %s exists, %v\n", email, err)
		return false
	}
	return count > 0
}
