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

/*
func TestList(t *testing.T) {
	c := NewClient("http://localhost:8080/v1/", time.Minute)
	var accounts account.Accounts
	c.List(account.Provider).As(&accounts)
	assert.Equal(t, 0, len(accounts))

	c.List(account.Provider, WithPageSize(1), WithPageNumber(1)).As(&accounts)
	assert.Equal(t, 0, len(accounts))
}

func TestList2(t *testing.T) {
	c := NewClient("http://localhost:8080/v1/", time.Minute)
	var accounts account.Accounts
	c.List(account.Provider).As(&accounts)
	assert.Equal(t, 1, len(accounts))

	c.List(account.Provider, WithPageSize(1), WithPageNumber(1)).As(&accounts)
	assert.Equal(t, 2, len(accounts))
}
*/
/*
func TestList(t *testing.T) {
	c := NewClient("http://localhost:8080/v1/", time.Minute)
	var accounts account.Accounts
	c.List(account.Provider).As(&accounts)
	assert.Equal(t, len(accounts), 3)

	c.List(account.Provider, WithPageSize(1), WithPageNumber(1)).As(&accounts)
	assert.Equal(t, len(accounts), 1)
}
*/
/*
func TestCreate(t *testing.T) {
	c := NewClient("http://localhost:8080/v1/", time.Minute)
	attributes := account.Attributes{
		Country:      "GB",
		BaseCurrency: "GBP",
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		BIC:          "NWBKGB22",
	}
	acc := account.NewAccount("", "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c", attributes)
	resp, err := c.Create(account.Provider, acc)

	assert.NoError(t, err)
	assert.False(t, resp.IsErrorStatus())

	var newAcc account.Account
	resp.As(&newAcc)
	fmt.Println(newAcc)
}*/

func TestIntegration(t *testing.T) {
	log.Println("Start of the integration test.")

	url := os.Getenv("API_URL")
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
}
