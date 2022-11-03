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
