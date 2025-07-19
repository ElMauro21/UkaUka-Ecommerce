package payu

import (
	"database/sql"

	"github.com/ElMauro21/UkaUkafb/helpers/cart"
)

func SaveTransactionItems(db *sql.DB, transactionID int, cartItems []cart.CartItem) error {
	for _, item := range cartItems {
		_, err := db.Exec(`
			INSERT INTO transaction_items (transaction_id, product_id, quantity)
			VALUES (?, ?, ?)`,
			transactionID, item.ProductID, item.Quantity,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
