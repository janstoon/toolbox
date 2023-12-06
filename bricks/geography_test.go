package bricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestLookupCountries(t *testing.T) {
	assert.Nil(t, bricks.LookupCountryByIsoAlphaTwoCode("XY"))
	assert.Nil(t, bricks.LookupCountryByTelephoneCode("990"))

	assert.NotNil(t, bricks.LookupCountryByIsoAlphaTwoCode(neverland.Codes.IsoAlphaTwo))
	assert.NotNil(t, bricks.LookupCountryByTelephoneCode(neverland.Codes.Telephone))
}
