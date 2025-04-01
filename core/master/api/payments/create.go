package paymentsapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/sellix"
	"api/core/models/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	Route.NewSub(server.NewRoute("/create", func(w http.ResponseWriter, r *http.Request) {
		type Status struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			ID      int    `json:"id"`
		}
		switch strings.ToLower(r.Method) {
		case "post":
			ok, user := sessions.IsLoggedIn(w, r)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
	
			Currency := r.PostFormValue("coin")
			//Coupon := r.PostFormValue("coupon")
			Amount, err := strconv.ParseFloat(r.PostFormValue("amount"), 64)
			if err != nil {
				json.NewEncoder(w).Encode(&Status{Status: "error", Message: err.Error(), ID: 0})
				return
			}
	
			// NowPayments CreatePayment call
			response, err := sellix.Manager.CreatePayment(Amount, models.Config.Autobuy.Flat, Currency, strconv.Itoa(user.ID))
			if err != nil {
				json.NewEncoder(w).Encode(&Status{Status: "error", Message: err.Error(), ID: 0})
				return
			}
	
			// Updated response handling
			id, err := database.Container.NewSale(&database.Sale{
				UniqID:       response.PaymentID,
				Amount:       int(Amount),
				Parent:       user.ID,
				Coin:         Currency,
				Status:       response.Status,
				Product:      fmt.Sprintf("Recharge for %s $%.2f", user.Username, Amount),
				Date:         time.Now().Unix(),
				Address:      response.PayAddress,
				CryptoAmount: response.PayAmount,
			})
			if err != nil {
				json.NewEncoder(w).Encode(&Status{Status: "error", Message: err.Error(), ID: 0})
				return
			}
	
			json.NewEncoder(w).Encode(&Status{Status: "success", Message: "Successfully created invoice " + response.PaymentID, ID: id})
		}
	}))
}
