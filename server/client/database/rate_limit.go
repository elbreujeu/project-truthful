package database

import (
	"database/sql"
	"log"
	"project_truthful/models"
	"time"
)

func GetRateLimit(ip string, db *sql.DB) (models.RateLimit, error) {
	var rateLimit models.RateLimit
	err := db.QueryRow("SELECT * FROM rate_limit WHERE ip_address = ?", ip).Scan(&rateLimit.IpAddress, &rateLimit.RequestCount, &rateLimit.LastRequestTime)
	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO rate_limit (ip_address, request_count, last_updated) VALUES (?, 0, NOW())", ip)
		if err != nil {
			log.Printf("Error inserting rate limit for ip %s, %v\n", ip, err)
			return models.RateLimit{}, err
		}
		return models.RateLimit{IpAddress: ip, RequestCount: 0, LastRequestTime: time.Now()}, nil
	} else if err != nil {
		log.Printf("Error getting rate limit for ip %s, %v\n", ip, err)
		return models.RateLimit{}, err
	}
	return rateLimit, nil
}

func ResetRateLimit(ip string, db *sql.DB) error {
	_, err := db.Exec("UPDATE rate_limit SET request_count = 0, last_updated = NOW() WHERE ip_address = ?", ip)
	if err != nil {
		log.Printf("Error resetting rate limit for ip %s, %v\n", ip, err)
		return err
	}
	return nil
}

func IncrementRateLimit(ip string, db *sql.DB) error {
	_, err := db.Exec("UPDATE rate_limit SET request_count = request_count + 1, last_updated = NOW() WHERE ip_address = ?", ip)
	if err != nil {
		log.Printf("Error incrementing rate limit for ip %s, %v\n", ip, err)
		return err
	}
	return nil
}
