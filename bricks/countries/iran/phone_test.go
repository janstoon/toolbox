package iran_test

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
}
