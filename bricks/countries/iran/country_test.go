package iran_test

import (
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestLookupCountries(t *testing.T) {
	ir := tricks.PtrVal(bricks.LookupCountryByIsoAlphaTwoCode("IR"))
	assert.Equal(t, "IR", ir.Codes.IsoAlphaTwo)
	assert.Equal(t, "IRN", ir.Codes.IsoAlphaThree)
	assert.Equal(t, "98", ir.Codes.Telephone)

	ir = tricks.PtrVal(bricks.LookupCountryByTelephoneCode("98"))
	assert.Equal(t, "IR", ir.Codes.IsoAlphaTwo)
	assert.Equal(t, "IRN", ir.Codes.IsoAlphaThree)
	assert.Equal(t, "98", ir.Codes.Telephone)
}
