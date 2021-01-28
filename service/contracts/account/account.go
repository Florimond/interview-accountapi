package account

import (
	"github.com/Florimond/interview-accountapi/service/contracts"
	"github.com/google/uuid"
)

// Provider returns the provider for the accounts
var Provider = contracts.NewStaticProvider("organisation/accounts")

// Attributes embed the attributes specific to an account
type Attributes struct {
	BankID       string `json:"bank_id"`
	BankIDCode   string `json:"bank_id_code"`
	BaseCurrency string `json:"base_currency"`
	BIC          string `json:"bic"`
	Country      string `json:"country"`
}

// Account is a struct containing all the attributes of a account record
type Account struct {
	contracts.RecordInfo
	Attributes Attributes `json:"attributes"`
}

// NewAccount creates a new Account, filling some necessary record info
func NewAccount(ID string, orgID string, attributes Attributes) *Account {
	if len(ID) == 0 {
		ID = uuid.New().String()
	}
	return &Account{
		RecordInfo: contracts.RecordInfo{
			ID:             ID,
			OrganisationID: orgID,
			Type:           "accounts",
		},
		Attributes: attributes,
	}
}

// Accounts represent a list of accounts.
// The List function will unmarshal to a simple slice of the Account struct, but it could be more
// complex in the future, hence the abstraction.
type Accounts []Account
