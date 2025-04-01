package sellix

import (
	"api/core/models"
	"io"
	"net/http"
	"strings"
)

var (
	Manager = NewClient(models.Config.Autobuy.Key)
)

type Client struct {
	client *http.Client
	APIKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		client: http.DefaultClient,
		APIKey: apiKey,
	}
}

func (c *Client) CreateRequest(method, data, endpoint string) (*http.Request, error) {
	r, err := http.NewRequest(method, RootURL+endpoint, func() io.Reader {
		if strings.ToLower(method) == "post" {
			return strings.NewReader(data)
		}
		return nil
	}())
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-api-key", models.Config.Autobuy.Key)
	return r, nil
}
