package plans

type Addon struct {
	Price   int 	`json:"price"`
	Fiat    string  `json:"fiat"`
	Upgrade string  `json:"upgrade"`
	Text	string  `json:"text"`
	Value   int     `json:"value"`
	Expiry  int     `json:"expiry"`
	Type	string	`json:"type"`
}
