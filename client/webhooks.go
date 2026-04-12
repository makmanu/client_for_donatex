package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/makmanu/client_for_donatex/config"
)

type Webhook struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	ClientId     string `json:"clientId"`
	EventType    string `json:"eventType"`
	IsActive     bool   `json:"isActive"`
	FailureCount int    `json:"failureCount"`
}

func (c *Client) CreateWebhook(secret *config.Secret) (*Webhook, error) {
	body := map[string]string{
		"url":       secret.WebhookURL,
		"secret":    secret.WebhookSecret,
		"eventType": "DonationCreated",
		"clientId":  secret.ClientId,
	}
	resp, err := c.DoRequest("POST", "webhooks/subscriptions", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Webhook creation failed. Status: %s, Body: %s\n", resp.Status, string(bodyBytes))
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
