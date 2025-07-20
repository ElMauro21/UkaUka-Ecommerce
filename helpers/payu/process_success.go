package payu

import (
	"database/sql"
	"fmt"
)

func ProcessSuccessfulTransaction(db *sql.DB, refCode string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var transactionID int
	var status string
	err = tx.QueryRow(`
		SELECT id, status FROM transactions WHERE reference_code = ?`, refCode).
		Scan(&transactionID, &status)
	if err != nil {
		return err
	}

	if status != "pending" {
		return fmt.Errorf("transaction already processed")
	}

	rows, err := tx.Query(`
		SELECT product_id, quantity FROM transaction_items WHERE transaction_id = ?`, transactionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, quantity int
		if err := rows.Scan(&productID, &quantity); err != nil {
			return err
		}

		res, err := tx.Exec(`
			UPDATE products SET quantity = quantity - ?
			WHERE id = ? AND quantity >= ?`, quantity, productID, quantity)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return fmt.Errorf("not enough stock for product %d", productID)
		}
	}

	_, err = tx.Exec(`
		UPDATE transactions SET status = 'completed' WHERE id = ?`, transactionID)
	if err != nil {
		return err
	}

	var userID sql.NullInt64
	var sessionID sql.NullString

	err = tx.QueryRow(`SELECT user_id, session_id FROM transactions WHERE id = ?`, transactionID).
		Scan(&userID, &sessionID)
	if err != nil {
		return err
	}

	if userID.Valid {
		_, err = tx.Exec(`
			DELETE FROM cart_items
			WHERE cart_id IN (
				SELECT id FROM carts WHERE user_id = ?
			)`, userID.Int64)
		if err != nil {
			return err
		}
	} else if sessionID.Valid {
		_, err = tx.Exec(`
			DELETE FROM cart_items
			WHERE cart_id IN (
				SELECT id FROM carts WHERE session_id = ?
			)`, sessionID.String)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}