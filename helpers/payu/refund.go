package payu

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
)

type RefundRequest struct {
	Language       string `json:"language"`
	Command        string `json:"command"`
	Test           bool   `json:"test"`
	Merchant       struct {
		ApiKey    string `json:"apiKey"`
		ApiLogin  string `json:"apiLogin"`
	} `json:"merchant"`
	Transaction struct {
		Order struct {
			Id string `json:"id"`
		} `json:"order"`
		ParentTransactionId string `json:"parentTransactionId"`
		Reason              string `json:"reason"`
		Type                string `json:"type"`
		PaymentMethod       string `json:"paymentMethod,omitempty"`
	} `json:"transaction"`
}

type RefundResponse struct {
	Code    string `json:"code"`
	Error   string `json:"error,omitempty"`
}

func RefundTransactionWithLog(db *sql.DB, transactionID int, payuTransactionID string) {
    status := "success"
    message := ""

    err := RefundTransaction(payuTransactionID)
    if err != nil {
        status = "fail"
        message = err.Error()
    }

    // Log de intento
    _, logErr := db.Exec(`
        INSERT INTO refund_attempts (transaction_id, payu_transaction_id, status, message)
        VALUES (?, ?, ?, ?)`,
        transactionID, payuTransactionID, status, message,
    )
    if logErr != nil {
        fmt.Printf("[Tx:%d] Failed to log refund attempt: %v\n", transactionID, logErr)
    }

    // 🔹 Si falla el reembolso
    if status == "fail" {
        fmt.Printf("[Tx:%d] Refund failed: %s\n", transactionID, message)

        var email, fullName string
        if err := db.QueryRow(`
            SELECT email, full_name FROM shipping_info WHERE transaction_id = ?`,
            transactionID).Scan(&email, &fullName); err != nil {
            fmt.Printf("[Tx:%d] Could not fetch email for refund failure: %v\n", transactionID, err)
            return
        }

        if email != "" {
            go func() {
                if err := auth.SendRefundFailureEmail(email, fullName); err != nil {
                    fmt.Printf("[Tx:%d] Failed to send refund failure email: %v\n", transactionID, err)
                } else {
                    fmt.Printf("[Tx:%d] Refund failure email sent\n", transactionID)
                }
            }()
        }
        return
    }

    // 🔹 Si el reembolso fue exitoso
    _, updateErr := db.Exec(`UPDATE transactions SET status = 'refunded' WHERE id = ?`, transactionID)
    if updateErr != nil {
        fmt.Printf("[Tx:%d] Failed to mark transaction as refunded: %v\n", transactionID, updateErr)
    }

    // Obtener datos de cliente
    var email, fullName string
    if err := db.QueryRow(`
        SELECT email, full_name FROM shipping_info WHERE transaction_id = ?`, transactionID).
        Scan(&email, &fullName); err != nil {
        fmt.Printf("[Tx:%d] Could not fetch email for refund confirmation: %v\n", transactionID, err)
        return
    }

    // Resumen de productos
    var productSummary string
    rows, err := db.Query(`
        SELECT p.name, ti.quantity
        FROM transaction_items ti
        JOIN products p ON p.id = ti.product_id
        WHERE ti.transaction_id = ?`, transactionID)
    if err == nil {
        defer rows.Close()
        var summary string
        for rows.Next() {
            var name string
            var qty int
            if err := rows.Scan(&name, &qty); err == nil {
                summary += fmt.Sprintf("%s x%d, ", name, qty)
            }
        }
        productSummary = summary
    }

    // Envío de correo en segundo plano
    if email != "" {
        go func() {
            if err := auth.SendRefundEmail(email, fullName, productSummary); err != nil {
                fmt.Printf("[Tx:%d] Failed to send refund email: %v\n", transactionID, err)
            } else {
                fmt.Printf("[Tx:%d] Refund email sent successfully\n", transactionID)
            }
        }()
    }

    fmt.Printf("[Tx:%d] Refund triggered successfully and transaction marked as refunded\n", transactionID)
}

func RefundTransaction(payuTransactionID string) error {
    env := os.Getenv("PAYU_ENV")
    apiURL := "https://sandbox.api.payulatam.com/payments-api/4.0/service.cgi"
    test := true

    if env == "production" {
        apiURL = "https://api.payulatam.com/payments-api/4.0/service.cgi"
        test = false
    }

    apiKey := os.Getenv("API_KEY")
    apiLogin := os.Getenv("API_LOGIN")
    if apiKey == "" || apiLogin == "" {
        return fmt.Errorf("PayU credentials missing")
    }

    req := RefundRequest{
        Language: "es",
        Command:  "void",
        Test:     test,
    }
    req.Merchant.ApiKey = apiKey
    req.Merchant.ApiLogin = apiLogin
    req.Transaction.ParentTransactionId = payuTransactionID
    req.Transaction.Reason = "Stock agotado"
    req.Transaction.Type = "VOID"

    payload, err := json.Marshal(req)
    if err != nil {
        return err
    }

    resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        return fmt.Errorf("refund failed with status: %s", resp.Status)
    }

    return nil
}