package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Florimond/interview-accountapi/service/contracts/account"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	log.Println("Start of the integration test.")

	url := os.Getenv("API_URL")
	if len(url) == 0 {
		url = "http://localhost:8080/v1"
	}
	client := NewClient(url, time.Minute)
	ctx := context.Background()

	log.Println("Create an account.")
	attributes := account.Attributes{
		Country:      "GB",
		BaseCurrency: "GBP",
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		BIC:          "NWBKGB22",
	}
	acc := account.NewAccount("", "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c", attributes)
	resp, err := client.Create(ctx, account.Provider, acc)
	require.NoError(t, err)

	log.Println("Account created successfully, unmarshal the response.")
	var newAcc account.Account
	err = resp.As(&newAcc)
	require.NoError(t, err)

	log.Println("Fetch the newly created account by ID.")
	resp, err = client.FindByID(ctx, account.Provider, newAcc.ID)
	require.NoError(t, err)

	log.Println("Account fetched successfully, unmarshal the response.")
	var fetchAcc account.Account
	err = resp.As(&fetchAcc)
	require.NoError(t, err)

	log.Println("List the accounts page 1.")
	resp, err = client.List(ctx, account.Provider, WithPageSize(1))
	require.NoError(t, err)

	log.Println("List retrieved successfully, unmarshal the response.")
	var accLst account.Accounts
	err = resp.As(&accLst)
	require.NoError(t, err)
	require.GreaterOrEqual(t, 1, len(accLst))
	//require.Equal(t, newAcc.ID, accLst[0].ID)

	log.Println("Delete the created account.")
	resp, err = client.Delete(ctx, account.Provider, newAcc.ID, newAcc.Version)
	require.NoError(t, err)

	log.Println("Attempt to fetch the deleted account.")
	resp, err = client.FindByID(ctx, account.Provider, newAcc.ID)
	require.Error(t, err)
	_, ok := IsHTTPError(err)
	require.True(t, ok)
}
