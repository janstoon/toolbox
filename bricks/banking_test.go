package bricks_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
)

func TestParseIban(t *testing.T) {
	bricks.RegisterBbanValidator(neverland.Codes.IsoAlphaTwo, func(bban string) error {
		if !strings.HasPrefix(bban, "014") {
			return bricks.ErrIbanInvalidBban
		}

		return nil
	})

	ibanStr := "NV670143022491272900023641"
	iban, err := bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.NoError(t, err)
	assert.NotNil(t, iban)
	assert.Equal(t, ibanStr, iban.String())
	assert.Equal(t, iban.Country, bricks.LookupCountryByIsoAlphaTwoCode(neverland.Codes.IsoAlphaTwo))

	ibanStr = "NV49"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanIncorrectLength)
	assert.Nil(t, iban)

	ibanStr = "NV670143022491272900023641000000000"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanIncorrectLength)
	assert.Nil(t, iban)

	ibanStr = "NV670143022491272900023642"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanCheckFailure)
	assert.Nil(t, iban)

	ibanStr = "NV660143022491272900023641"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanCheckFailure)
	assert.Nil(t, iban)

	ibanStr = "XY870143022491272900023730"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanUnknownCountry)
	assert.Nil(t, iban)

	ibanStr = "NV860113022491272900023641"
	iban, err = bricks.ParseInternationalBankAccountNumber(ibanStr)
	require.ErrorIs(t, err, bricks.ErrIbanInvalidBban)
	assert.Nil(t, iban)
}

func TestParsePan(t *testing.T) {
	panStr := "6274129005473742"
	pan, err := bricks.ParsePrimaryAccountNumber(panStr)
	require.NoError(t, err)
	assert.NotNil(t, pan)
	assert.Equal(t, panStr, pan.String())
}
