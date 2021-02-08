package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/Florimond/interview-accountapi/service/contracts/account"
	"github.com/stretchr/testify/assert"
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
}
