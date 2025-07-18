package payu

import "database/sql"

func SaveShippingInfo(db *sql.DB,fullName, idNumber, phone, email, state, city, neighborhood, address string, transactionID int)error{

	_, err := db.Exec(`
        INSERT INTO shipping_info (
            transaction_id, full_name, id_number, phone, email, state, city, neighborhood, address
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        transactionID,
        fullName,
        idNumber,
        phone,
        email,
        state,
        city,
        neighborhood,
        address,
    )

	return err
}
