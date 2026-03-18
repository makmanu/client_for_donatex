package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL string
	APIKey  string
}

type Donation struct {
	ID                string  `json:"id"`
	Username          string  `json:"username"`
	Message           string  `json:"message"`
	Currency          string  `json:"currency"`
	Amount            float64 `json:"amount"`
	AmountInRub       float64 `json:"amountInRub"`
	Timestamp         string  `json:"timestamp"`
	WithAiResponse    bool    `json:"withAiResponse"`
	AiResponse        string  `json:"aiResponse"`
	IsTest            bool    `json:"isTest"`
	IsPotentiallyUnsafe bool  `json:"isPotentiallyUnsafe"`
	WasShown          bool    `json:"wasShown"`
	IsFeePaidByUser   bool    `json:"isFeePaidByUser"`
	VoiceFilePath     string  `json:"voiceFilePath"`
	PaidVoice         string `json:"paidVoice"`
	MusicLink         string  `json:"musicLink"`
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (c *Client) DoRequest(method, endpoint string, body interface{}, queryParams map[string]string) (*http.Response, error) {
	// Parse the base URL
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Parse the endpoint
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %v", err)
	}

	// Combine base URL and endpoint
	fullURL := baseURL.ResolveReference(endpointURL)

	// Add query parameters
	query := fullURL.Query()
	for key, value := range queryParams {
		query.Set(key, value)
	}
	if c.APIKey != "" {
		query.Set("token", c.APIKey)
	}
	fullURL.RawQuery = query.Encode()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fullURL.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Example method for getting donations
func (c *Client) GetDonations(skip, take int, hideTest string) error {
	queryParams := map[string]string{
		"skip": fmt.Sprintf("%d", skip),
		"take": fmt.Sprintf("%d", take),
		"hideTest": hideTest,
	}

	resp, err := c.DoRequest("GET", "donations", nil, queryParams)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var donations []Donation
	err = json.Unmarshal(body, &donations)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	fmt.Printf("Retrieved %d donations (skip: %d, take: %d):\n", len(donations), skip, take)
	for i, donation := range donations {
		fmt.Printf("%d. Username: %s\n\t Message: %s,\n\t Amount: %.2f %s (%.2f RUB),\n\t Timestamp: %s\n\n",
			i+1, donation.Username, donation.Message, donation.Amount, donation.Currency, donation.AmountInRub, donation.Timestamp)
	}

	return nil
}

func (c *Client) TestDonations(amount float32, username, message, currency string, withAiResponse bool) error {
	body := map[string]interface{}{
		"amount":          amount,
		"username":        username,
		"message":         message,
		"currency":        currency,
		"withAiResponse":  withAiResponse,
	}

	resp, err := c.DoRequest("POST", "test-donation", body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response body: %s\n", string(bodyBytes))
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return nil
}