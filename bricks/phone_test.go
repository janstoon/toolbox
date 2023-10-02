package bricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestParsePhoneNumber(t *testing.T) {
	pn, err := bricks.ParsePhoneNumber("00989123456789")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.True(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+989123456789")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.True(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989123456789")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.True(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+98 912 345 67 89")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.True(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("00982122334455")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.False(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+982122334455")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.False(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("982122334455")
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode("IR"), pn.Country)
	assert.False(t, pn.DefaultOperator.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+9718005625926")
	assert.ErrorIs(t, err, bricks.ErrPhoneNumberUnknownCountry)
	assert.Nil(t, pn)

	pn, err = bricks.ParsePhoneNumber("+97125236227")
	assert.ErrorIs(t, err, bricks.ErrPhoneNumberUnknownCountry)
	assert.Nil(t, pn)

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
		pn, err = bricks.ParsePhoneNumber(number)
		assert.ErrorIs(t, err, bricks.ErrInvalidInput)
		assert.Nil(t, pn)
	}
}
