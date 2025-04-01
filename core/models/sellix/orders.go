package sellix

import "encoding/json"

type OrdersResp struct {
	Status int `json:"status"`
	Data   struct {
		Orders []struct {
			Uniqid  string `json:"uniqid"`
			Product struct {
				Currency string `json:"currency"`
				Price    int    `json:"price"`
			} `json:"product"`
		} `json:"orders"`
	} `json:"data"`
}

func (c *Client) GetPaymentStatus(paymentID string) (*PaymentResp, error) {
	req, err := c.CreateRequest("GET", "", "/payment/"+paymentID)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pResp PaymentResp
	if err := json.NewDecoder(resp.Body).Decode(&pResp); err != nil {
		return nil, err
	}

	return &pResp, nil
}