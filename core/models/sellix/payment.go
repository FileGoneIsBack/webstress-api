package sellix

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"api/core/models"
)

type PaymentReq struct {
	PriceAmount   float64 `json:"price_amount"`
	PriceCurrency string  `json:"price_currency"`
	PayCurrency   string  `json:"pay_currency"`
	OrderID       string  `json:"order_id"`
	IPNCallback   string  `json:"ipn_callback_url"`
}

type PaymentResp struct {
	Status    string  `json:"payment_status"`
	PaymentID string  `json:"payment_id"`
	PayAmount float64 `json:"pay_amount"`
	PayAddress string `json:"pay_address"`
}

func (c *Client) CreatePayment(amount float64, fiat, crypto, orderID string) (*PaymentResp, error) {
	payment := &PaymentReq{
		PriceAmount:   amount,
		PriceCurrency: fiat,
		PayCurrency:   crypto,
		OrderID:       orderID,
		IPNCallback:   models.Config.Domain + "api/payments/webhook",
	}

	body, err := json.Marshal(payment)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal payment request: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Sending payment request: %s", string(body))

	req, err := c.CreateRequest("POST", string(body), Payments)
	if err != nil {
		log.Printf("[ERROR] Failed to create payment request: %v", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed to send payment request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[INFO] Received response with status: %s", resp.Status)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Response Body: %s", string(respBody))

	var pResp PaymentResp
	if err := json.Unmarshal(respBody, &pResp); err != nil {
		log.Printf("[ERROR] Failed to decode payment response: %v", err)
		return nil, err
	}

	if pResp.Status == "" {
		log.Println("[ERROR] Invalid payment response: missing status field")
		return nil, errors.New("invalid payment response, min deposit is 25$")
	}

	log.Printf("[INFO] Payment created successfully: %+v", pResp)
	return &pResp, nil
}
