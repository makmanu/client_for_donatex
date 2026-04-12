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

func (c *Client) CreateWebhook() (*Webhook, error) {
	
	secret, err := config.LoadSecret("secret.yaml")
	if err != nil {
		return nil, fmt.Errorf("Failed to load secret: %v", err)
	}

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

func (c *Client) GetWebhooks() ([]Webhook, error) {
	resp, err := c.DoRequest("GET", "webhooks/subscriptions", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var webhooks []Webhook
	err = json.Unmarshal(body, &webhooks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return webhooks, nil
}

func (c *Client) DeleteWebhook(webhookId string) error {
	resp, err := c.DoRequest("DELETE", fmt.Sprintf("webhooks/subscriptions/%s", webhookId), nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return nil
}