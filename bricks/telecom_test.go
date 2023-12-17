package bricks_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
)

func TestNetworkOperatorsByCountryIsoAlphaTwoCode(t *testing.T) {
	bricks.RegisterNetworkOperators(neverland.Codes.IsoAlphaTwo, bricks.NetworkOperator{
		Name:    "NeverTel",
		Virtual: false,
	})
	bricks.RegisterNetworkOperators(neverland.Codes.IsoAlphaTwo,
		bricks.NetworkOperator{
			Name:    "ZeroTel",
			Virtual: false,
		},
		bricks.NetworkOperator{
			Name:    "NilTel",
			Virtual: true,
		},
	)

	nn := bricks.NetworkOperatorsByCountryCode(neverland.Codes.IsoAlphaTwo)
	assert.Len(t, nn, 3)
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "NeverTel",
		Virtual: false,
	})
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "ZeroTel",
		Virtual: false,
	})
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "NilTel",
		Virtual: true,
	})

	assert.PanicsWithValue(t, bricks.ErrNotFound, func() {
		bricks.RegisterNetworkOperators("XY", bricks.NetworkOperator{
			Name:    "PanicTel",
			Virtual: false,
		})
	})
	assert.Nil(t, bricks.NetworkOperatorsByCountryCode("XY"))
}

func TestParsePhoneNumber(t *testing.T) {
	bricks.RegisterPhoneNumberResolver(
		neverland.Codes.Telephone,
		func(localNumber string) (*bricks.PhoneNumberMetadata, error) {
			if strings.HasPrefix(localNumber, "123") {
				return &bricks.PhoneNumberMetadata{
					Mobile: true,
					Operator: bricks.NetworkOperator{
						Name:    "NeverTel",
						Virtual: false,
					},
				}, nil
			}

			return nil, errors.Join(bricks.ErrInvalidInput, bricks.ErrUnknownNetworkOperator)
		},
	)

	validNumbers := []string{
		"00999123456789",
		"+999123456789",
		"999123456789",
		"+999 123 45 67 89",
		"999 123-456-789",
		"+999(123) 456-789",
		"999(123) 456-789",
		"999(123) 456 789",
	}
	for _, number := range validNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		require.NoError(t, err)
		assert.NotNil(t, pn)
		assert.Equal(t, "+999123456789", pn.String())
		assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(neverland.Codes.IsoAlphaTwo), pn.Country)
		assert.True(t, pn.Mobile)
		assert.False(t, pn.DefaultOperator.Virtual)
		assert.Equal(t, "NeverTel", pn.DefaultOperator.Name)

		assert.Equal(t, tricks.PtrVal(pn), bricks.MustParsePhoneNumber(number))
	}

	unregisteredNumbers := []string{
		"+9718005625926",
		"+97125236227",
		"(123) 456-789",
		"123456789",
	}
	for _, number := range unregisteredNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		require.ErrorIs(t, err, bricks.ErrPhoneNumberUnknownCountry)
		assert.Nil(t, pn)

		require.Panics(t, func() {
			bricks.MustParsePhoneNumber(number)
		})
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
		"+999321456789",
	}
	for _, number := range invalidNumbers {
		pn, err := bricks.ParsePhoneNumber(number)
		require.ErrorIs(t, err, bricks.ErrInvalidInput)
		assert.Nil(t, pn)

		require.Panics(t, func() {
			bricks.MustParsePhoneNumber(number)
		})
	}
}
