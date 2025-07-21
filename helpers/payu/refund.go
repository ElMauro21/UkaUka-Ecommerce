package payu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
