package payu

import (
	"database/sql"
	"fmt"
	"time"
)

func CreateTransaction(db *sql.DB, userID *int, sessionID *string, totalAmount float64) (string,error) {
	now := time.Now()
	yearMonth := now.Format("200601")

	var count int 
	err := db.QueryRow(`
	SELECT COUNT(*) FROM transactions 
	WHERE strftime('%Y%m',created_at) = ?`,yearMonth).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("counting transactions: %v", err)
	}

	refCode := fmt.Sprintf("UKA%s%02d",yearMonth,count+1)

	_,err = db.Exec(`
	INSERT INTO transactions (user_id, session_id, reference_code, total_amount) 
	VALUES (?, ?, ?, ?)`,userID,sessionID, refCode,totalAmount)
	if err != nil {
		return "", fmt.Errorf("inserting transaction: %v", err)
	}

	return refCode, nil

}	