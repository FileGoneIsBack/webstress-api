package paymentsapi

import (
	"api/core/database/site"
	"api/core/master/sessions"
	"api/core/models/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	// Mapping short coin names to full names
	var coinNameMap = map[string]string{
		"btc": "bitcoin",
		"eth": "ethereum",
	}

	Route.NewSub(server.NewRoute("/status", func(w http.ResponseWriter, r *http.Request) {
		type Payment struct {
			ID            int     `json:"id"`
			Amount        int     `json:"amount"`
			CryptoAmount  float64 `json:"crypto_amount"`
			CryptoAddress string  `json:"crypto_address"`
			CryptoCoin    string  `json:"crypto_coin"`
			QrCode        string  `json:"qr_code"`
			Recieved      float64 `json:"recieved"`
			Status        string  `json:"status"`
			Date          int64   `json:"date"`
		}
		type Status struct {
			Status  string   `json:"status"`
			Message string   `json:"message"`
			Payment *Payment `json:"payment_info"`
		}
		switch strings.ToLower(r.Method) {
		case "post":
			ok, _ := sessions.IsLoggedIn(w, r)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			id, _ := strconv.Atoi(r.PostFormValue("payment_id"))
			payment, err := site.Container.GetSale(id)
			if err != nil {
				json.NewEncoder(w).Encode(&Status{Status: "error", Message: err.Error(), Payment: &Payment{}})
				return
			}
			
			// Get full coin name
			fullCoinName, exists := coinNameMap[strings.ToLower(payment.Coin)]
			if !exists {
				fullCoinName = payment.Coin // fallback to original coin name if it's not in the map
			}

			// Generate the QR code URL with the full coin name
			qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=280x280&data=%s:%s?amount=%s", fullCoinName, payment.Address, fmt.Sprint(payment.CryptoAmount))
			// Send the response with the updated QR code URL
			json.NewEncoder(w).Encode(&Status{
				Status:  "success",
				Message: "transaction information",
				Payment: &Payment{
					ID:            payment.ID,
					Amount:        payment.Amount,
					CryptoAmount:  payment.CryptoAmount,
					CryptoAddress: payment.Address,
					CryptoCoin:    payment.Coin,
					QrCode:        qrCodeURL,
					Recieved:      payment.Recieved,
					Status:        payment.Status,
					Date:          payment.Date,
				},
			})
		}
	}))
}
