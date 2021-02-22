package account

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestNewAccount(t *testing.T) {
	accountID := uuid.New().String()
	orgID := uuid.New().String()
	attributes := Attributes{
		BankID: uuid.New().String(),
	}
	account1 := NewAccount(accountID, orgID, attributes)
	account2 := NewAccount("", orgID, attributes)

	// Test record creation with provided ID.
	assert.NotNil(t, account1)
	assert.Equal(t, accountID, account1.ID)
	assert.Equal(t, orgID, account1.OrganisationID)
	assert.Equal(t, attributes.BankID, account1.Attributes.BankID)

	// Test record created without provided ID.
	assert.NotNil(t, account2)
	assert.NotEmpty(t, account2.ID)
}

func TestPath(t *testing.T) {
	assert.Equal(t, "organisation/accounts/", Path(""))
	assert.Equal(t, "organisation/accounts/123", Path("123"))
}
