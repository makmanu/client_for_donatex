package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Webhook struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	ClientId     string `json:"clientId"`
	EventType    string `json:"eventType"`
	IsActive     bool   `json:"isActive"`
	FailureCount int    `json:"failureCount"`
}

func (c *Client) CreateWebhook(webhookURL, eventType, clientId, secret string) (*Webhook, error) {
	body := map[string]string{
		"url":       webhookURL,
		"eventType": eventType,
		"clientId":  clientId,
		"secret":    secret,
	}
	resp, err := c.DoRequest("POST", "/webhooks", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	err = json.Unmarshal(bodyBytes, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}
