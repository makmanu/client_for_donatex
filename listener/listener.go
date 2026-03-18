package listener

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
)

func StartListener(cfg *config.Config) {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received webhook request")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var donation client.Donation
		err := json.NewDecoder(r.Body).Decode(&donation)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Output to "chat" (console)
		fmt.Printf("💰 New donation from %s: %.2f %s - %s\n", donation.Username, donation.Amount, donation.Currency, donation.Message)

		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Starting listener on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
