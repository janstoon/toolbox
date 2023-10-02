package iran_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestParseBban(t *testing.T) {
	// correct iban - correct bban
	ibanStr := "IR490143022491272900023730"
	iban, err := bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.NoError(t, err)
	assert.NotNil(t, iban)
	assert.Equal(t, ibanStr, iban.String())
	assert.Equal(t, iban.Country, bricks.LookupCountryByIsoAlphaTwoCode("IR"))

	// correct iban - bban length = 21
	ibanStr = "IR07014582656540804450248"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanInvalidBban)
	assert.ErrorContains(t, err, "bban length incorrect")
	assert.Nil(t, iban)

	// correct iban - bban length = 23
	ibanStr = "IR5301458265654080445024885"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanInvalidBban)
	assert.ErrorContains(t, err, "bban length incorrect")
	assert.Nil(t, iban)

	// correct iban - invalid bi
	ibanStr = "IR670095826565408044502441"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanInvalidBban)
	assert.ErrorContains(t, err, "invalid bank identifier")
	assert.Nil(t, iban)
}
