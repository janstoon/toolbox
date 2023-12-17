package iran_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestNetworkOperatorsList(t *testing.T) {
	oo := bricks.NetworkOperatorsByCountryCode(iran.Iran.Codes.IsoAlphaTwo)
	assert.Len(t, oo, 13)
	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "MCI",
		Virtual: false,
	})
	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "MTN",
		Virtual: false,
	})
	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "Rightel",
		Virtual: false,
	})
	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "Taliya",
		Virtual: false,
	})

	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "TCI",
		Virtual: false,
	})

	assert.Contains(t, oo, bricks.NetworkOperator{
		Name:    "Shatel",
		Virtual: true,
	})
}

func TestParsePhoneNumber(t *testing.T) {
	pn, err := bricks.ParsePhoneNumber("00989123456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+989123456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989123456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+98 912 345 67 89")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989193456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989193456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.True(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989353456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989353456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.True(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MTN", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("00982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.Prepaid)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989123456")
	require.ErrorIs(t, err, bricks.ErrInvalidInput)
	require.ErrorContains(t, err, "local number length incorrect")
	assert.Nil(t, pn)

	pn, err = bricks.ParsePhoneNumber("9891234567890")
	require.ErrorIs(t, err, bricks.ErrInvalidInput)
	require.ErrorContains(t, err, "local number length incorrect")
	assert.Nil(t, pn)

	pn, err = bricks.ParsePhoneNumber("981023456789")
	require.ErrorIs(t, err, bricks.ErrInvalidInput)
	require.ErrorIs(t, err, bricks.ErrUnknownNetworkOperator)
	assert.Nil(t, pn)
}
