package bricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestParseIban(t *testing.T) {
	iban, err := bricks.ParseInternationalBankAccountNumber("IR490143022491272900023730")
	assert.NoError(t, err)
	assert.NotNil(t, iban)
}

func TestParsePan(t *testing.T) {

}
