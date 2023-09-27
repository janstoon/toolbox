package bricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestParseIban(t *testing.T) {
	ibanStr := "IR490143022491272900023730"
	iban, err := bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.NoError(t, err)
	assert.NotNil(t, iban)
	assert.Equal(t, ibanStr, iban.String())
	assert.Equal(t, iban.Country, bricks.LookupCountryByIsoAlphaTwoCode("IR"))

	ibanStr = "IR49"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanIncorrectLength)
	assert.Nil(t, iban)

	ibanStr = "IR490143022491272900023730000000000"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanIncorrectLength)
	assert.Nil(t, iban)

	ibanStr = "XY870143022491272900023730"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanUnknownCountry)
	assert.Nil(t, iban)

	ibanStr = "IR490143022491272900023731"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanCheckFailure)
	assert.Nil(t, iban)

	ibanStr = "IR480143022491272900023730"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	assert.ErrorIs(t, err, bricks.ErrIbanCheckFailure)
	assert.Nil(t, iban)
}

func TestParsePan(t *testing.T) {
	panStr := "6274129005473742"
	pan, err := bricks.ParsePrimaryAccountNumber(panStr)
	assert.NoError(t, err)
	assert.NotNil(t, pan)
	assert.Equal(t, panStr, pan.String())
}
