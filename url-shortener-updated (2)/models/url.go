package models

import (
	"database/sql"
	"time"
	"url-shortener/db"
)

type URL struct {
	ID          int
	UserID      int
	OriginalURL string
	ShortCode   string
	Expiry      time.Time
	CreatedAt   time.Time
}

func SaveURL(userID int, originalURL, shortCode string, expiry time.Time) error {
	_, err := db.DB.Exec(
		"INSERT INTO urls (user_id, original_url, short_code, expiry) VALUES (?, ?, ?, ?)",
		userID, originalURL, shortCode, expiry,
	)
	return err
}

func GetURLByCode(shortCode string) (*URL, error) {
	row := db.DB.QueryRow("SELECT id, user_id, original_url, short_code, expiry, created_at FROM urls WHERE short_code = ?", shortCode)

	var url URL
	err := row.Scan(&url.ID, &url.UserID, &url.OriginalURL, &url.ShortCode, &url.Expiry, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func GetUserURLs(userID int) ([]URL, error) {
	rows, err := db.DB.Query("SELECT id, original_url, short_code, expiry, created_at FROM urls WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		url.UserID = userID
		err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.Expiry, &url.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}


// DeleteExpiredURLs removes all URLs from the database that have expired.
func DeleteExpiredURLs(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM urls WHERE expiry < ?", time.Now())
	return err
}
