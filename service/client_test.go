package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Florimond/interview-accountapi/service/contracts/account"
)

func TestFormatOptions(t *testing.T) {
	tests := []struct {
		options  []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"a"}, "?a"},
		{[]string{"a", "b"}, "?a&b"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, formatOptions(test.options))
	}
}

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
	resp, err := client.Create(ctx, account.Path(""), acc)
	require.NoError(t, err)

	log.Println("Account created successfully, unmarshal the response.")
	var newAcc account.Account
	err = resp.As(&newAcc)
	require.NoError(t, err)

	log.Println("Fetch the newly created account by ID.")
	resp, err = client.FindByID(ctx, account.Path(newAcc.ID))
	require.NoError(t, err)

	log.Println("Account fetched successfully, unmarshal the response.")
	var fetchAcc account.Account
	err = resp.As(&fetchAcc)
	require.NoError(t, err)

	log.Println("List the accounts page 1.")
	resp, err = client.List(ctx, account.Path(""), WithPageSize(1))
	require.NoError(t, err)

	log.Println("List retrieved successfully, unmarshal the response.")
	var accLst account.Accounts
	err = resp.As(&accLst)
	require.NoError(t, err)
	require.GreaterOrEqual(t, 1, len(accLst))

	log.Println("Delete the created account.")
	resp, err = client.Delete(ctx, account.Path(newAcc.ID), newAcc.Version)
	require.NoError(t, err)

	log.Println("Attempt to fetch the deleted account.")
	resp, err = client.FindByID(ctx, account.Path(newAcc.ID))
	require.Error(t, err)
	_, ok := IsHTTPError(err)
	require.True(t, ok)
}
