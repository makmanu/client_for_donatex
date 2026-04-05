package listener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
)

func StartListener(cfg *config.Config) {
	logFile := cfg.LogFile
	if logFile == "" {
		logFile = "webhook.log"
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received webhook request")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))

		if err := logRequestToFile(logFile, r, body); err != nil {
			fmt.Printf("Failed to write webhook log: %v\n", err)
		}

		var payload client.WebhookPayload
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		donation := payload.Data

		fmt.Printf("💰 New donation from %s: %.2f %s - %s\n", donation.Username, donation.Amount, donation.Currency, donation.Message)

		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Starting HTTPS listener on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", cfg.Port), cfg.CertFile, cfg.KeyFile, nil))
}

func logRequestToFile(path string, r *http.Request, body []byte) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	header := fmt.Sprintf("[%s] %s %s %s from %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.String(), r.Proto, r.RemoteAddr)
	if _, err := file.WriteString(header); err != nil {
		return err
	}

	for name, values := range r.Header {
		if _, err := file.WriteString(fmt.Sprintf("%s: %s\n", name, strings.Join(values, ", "))); err != nil {
			return err
		}
	}

	if _, err := file.WriteString(fmt.Sprintf("\nBody:\n%s\n", string(body))); err != nil {
		return err
	}

	if _, err := file.WriteString("----\n"); err != nil {
		return err
	}

	return nil
}
