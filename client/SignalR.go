package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/makmanu/client_for_donatex/planner"
	"github.com/makmanu/client_for_donatex/structs"
	"github.com/philippseith/signalr"
)

type DonationReceiver struct {
	signalr.Hub
}

func (r *DonationReceiver) DonationCreated(donation structs.Donation) {
	fmt.Printf("💰 New donation from %s: %.2f %s - %s\n", donation.Username, donation.Amount, donation.Currency, donation.Message)
	planner.HandleDonation(donation)
}

func NewSignalRLogger() func(signalr.Party) error {
	f, err := os.OpenFile("signalr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger := kitlog.NewLogfmtLogger(f)

	// ВАЖНО: это уже OPTION, НЕ StructuredLogger
	return signalr.Logger(logger, true)
}

func ConnectWithTokenAutoReconnect(ctx context.Context, baseURL string, externalToken string) (signalr.Client, error) {
	u, err := url.Parse(baseURL + "/api/public-donations-hub")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("access_token", externalToken)
	u.RawQuery = q.Encode()

	connFactory := func() (signalr.Connection, error) {
		creationCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		return signalr.NewHTTPConnection(creationCtx, u.String())
	}

	signalrClient, err := signalr.NewClient(ctx,
		signalr.WithConnector(connFactory),
		signalr.WithReceiver(&DonationReceiver{}),
		signalr.KeepAliveInterval(15 * time.Second),
		NewSignalRLogger(),
	)
	if err != nil {
		return nil, err
	}

	signalrClient.Start()

	errCh := signalrClient.WaitForState(ctx, signalr.ClientConnected)
	if err := <-errCh; err != nil {
		return nil, err
	}

	log.Println("Connected with auto-reconnect enabled")

	return signalrClient, nil
}