package iran_test

import (
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestLookupCountries(t *testing.T) {
	ir := tricks.PtrVal(bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo))
	assert.Equal(t, iran.Iran.Codes.IsoAlphaTwo, ir.Codes.IsoAlphaTwo)
	assert.Equal(t, iran.Iran.Codes.IsoAlphaThree, ir.Codes.IsoAlphaThree)
	assert.Equal(t, iran.Iran.Codes.IsoNumeric, ir.Codes.IsoNumeric)
	assert.Equal(t, iran.Iran.Codes.IocAlphaThree, ir.Codes.IocAlphaThree)
	assert.Equal(t, iran.Iran.Codes.Telephone, ir.Codes.Telephone)

	ir = tricks.PtrVal(bricks.LookupCountryByTelephoneCode(iran.Iran.Codes.Telephone))
	assert.Equal(t, iran.Iran.Codes.IsoAlphaTwo, ir.Codes.IsoAlphaTwo)
	assert.Equal(t, iran.Iran.Codes.IsoAlphaThree, ir.Codes.IsoAlphaThree)
	assert.Equal(t, iran.Iran.Codes.IsoNumeric, ir.Codes.IsoNumeric)
	assert.Equal(t, iran.Iran.Codes.IocAlphaThree, ir.Codes.IocAlphaThree)
	assert.Equal(t, iran.Iran.Codes.Telephone, ir.Codes.Telephone)
}
