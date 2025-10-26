package jobs

import (
	"database/sql"
	"log"
	"time"
)

func JobTokenCleanup(db *sql.DB){
	go func() {
		for{
			time.Sleep(1 * time.Hour)

			result, err := db.Exec("DELETE FROM password_resets WHERE expires_at < ?",time.Now())
            if err != nil {
                log.Println("Failed to clean up expired tokens:", err)
                continue
            }

			rows, _ := result.RowsAffected()
			log.Printf("Cleaned up %d expired tokens\n", rows)
		}
	}()
}