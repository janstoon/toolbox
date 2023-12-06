package iran_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestNetworkOperatorsList(t *testing.T) {
	nn := bricks.NetworkOperatorsByCountryCode(iran.Iran.Codes.IsoAlphaTwo)
	assert.Len(t, nn, 13)
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "MCI",
		Virtual: false,
	})
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "MTN",
		Virtual: false,
	})
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "Rightel",
		Virtual: false,
	})
	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "Taliya",
		Virtual: false,
	})

	assert.Contains(t, nn, bricks.NetworkOperator{
		Name:    "TCI",
		Virtual: false,
	})

	assert.Contains(t, nn, bricks.NetworkOperator{
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
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+989123456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("989123456789")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+98 912 345 67 89")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+989123456789", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.True(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "MCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("00982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("+982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)

	pn, err = bricks.ParsePhoneNumber("982122334455")
	require.NoError(t, err)
	assert.NotNil(t, pn)
	assert.Equal(t, "+982122334455", pn.String())
	assert.Equal(t, bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo), pn.Country)
	assert.False(t, pn.Mobile)
	assert.False(t, pn.DefaultOperator.Virtual)
	assert.Equal(t, "TCI", pn.DefaultOperator.Name)
}
