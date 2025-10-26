package jobs

import (
	"database/sql"
	"log"
	"time"
)

func JobGuestCartCleanup(db *sql.DB) {
	go func() {
		for {
			time.Sleep(1 * time.Hour)

			threshold := time.Now().Add(-24 * time.Hour)

			result, err := db.Exec(`
				DELETE FROM carts
				WHERE user_id IS NULL
				AND created_at < ?
			`, threshold)

			if err != nil {
				log.Println("Error limpiando carritos de invitados:", err)
				continue
			}

			rows, _ := result.RowsAffected()
			log.Printf("Se eliminaron %d carritos de invitados vacíos y antiguos\n", rows)
		}
	}()
}