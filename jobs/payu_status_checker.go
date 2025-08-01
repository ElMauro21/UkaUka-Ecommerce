package jobs

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type StatusRequest struct {
	Language string `json:"language"`
	Command  string `json:"command"`
	Test     bool   `json:"test"`
	Merchant struct {
		ApiKey   string `json:"apiKey"`
		ApiLogin string `json:"apiLogin"`
	} `json:"merchant"`
	Details struct {
		TransactionID string `json:"transactionId"`
	} `json:"details"`
}

type StatusResponse struct {
	Code        string `json:"code"`
	Error       string `json:"error,omitempty"`
	Transaction struct {
		State string `json:"state"` 
	} `json:"transactionResponse"`
}

func CheckPendingTransactions(db *sql.DB) {
	go func() {
		for {

			fmt.Println("🔄 Checking pending transactions...")

			rows, err := db.Query(`SELECT id, payu_transaction_id FROM transactions WHERE status = 'pending' AND payu_transaction_id IS NOT NULL`)
			if err != nil {
				fmt.Printf("Error querying pending transactions: %v\n", err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var payuTxID string
				if err := rows.Scan(&id, &payuTxID); err != nil {
					fmt.Println("Scan error:", err)
					continue
				}

				state, err := GetPayUTransactionStatus(payuTxID)
				if err != nil {
					fmt.Printf("❌ Failed to check transaction %d: %v\n", id, err)
					continue
				}

				switch state {
				case "APPROVED":
					_, _ = db.Exec(`UPDATE transactions SET status = 'completed' WHERE id = ?`, id)
				case "DECLINED", "EXPIRED":
					_, _ = db.Exec(`UPDATE transactions SET status = 'failed' WHERE id = ?`, id)
				case "PENDING":
				default:
					fmt.Printf("Unknown state for transaction %d: %s\n", id, state)
				}
			}
		}
	}()
}

func GetPayUTransactionStatus(payuTxID string) (string, error) {
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
		return "", fmt.Errorf("missing PayU credentials")
	}

	req := StatusRequest{
		Language: "es",
		Command:  "TRANSACTION_RESPONSE_DETAIL",
		Test:     test,
	}
	req.Merchant.ApiKey = apiKey
	req.Merchant.ApiLogin = apiLogin
	req.Details.TransactionID = payuTxID

	body, _ := json.Marshal(req)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("PayU status request failed: %s", resp.Status)
	}

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return "", err
	}

	if statusResp.Code != "SUCCESS" {
		return "", fmt.Errorf("PayU status error: %s", statusResp.Error)
	}

	return statusResp.Transaction.State, nil
}
