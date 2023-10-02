package bricks_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestParsePhoneNumber(t *testing.T) {
	bricks.RegisterPhoneNumberLocator(neverland.Codes.Telephone, func(localNumber string) (*bricks.NetworkOperator, error) {
		if strings.HasPrefix(localNumber, "123") {
			return &bricks.NetworkOperator{
				Name:    "NeverTel",
				Mobile:  true,
				Virtual: false,
			}, nil
		}

		return nil, errors.Join(bricks.ErrInvalidInput, bricks.ErrUnknownNetworkOperator)
	})

	validNumbers := []string{
		"00999123456789",
		"+999123456789",
		"999123456789",
		"+999 123 45 67 89",
	}
	for _, number := range validNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		assert.NoError(t, err)
		assert.NotNil(t, pn)
		assert.Equal(t, "+999123456789", pn.String())
		assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(neverland.Codes.IsoAlphaTwo), pn.Country)
		assert.True(t, pn.DefaultOperator.Mobile)
		assert.False(t, pn.DefaultOperator.Virtual)
		assert.Equal(t, "NeverTel", pn.DefaultOperator.Name)
	}

	unregisteredNumbers := []string{
		"+9718005625926",
		"+97125236227",
	}
	for _, number := range unregisteredNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		assert.ErrorIs(t, err, bricks.ErrPhoneNumberUnknownCountry)
		assert.Nil(t, pn)
	}

	// invalid numbers
	invalidNumbers := []string{
		"",
		"1",
		"110",
		"119",
		"125",
		"911",
		"2233445",
		"22334455",
		"9123456789",
		"2122334455",
	}
	for _, number := range invalidNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		assert.ErrorIs(t, err, bricks.ErrInvalidInput)
		assert.Nil(t, pn)
	}
}