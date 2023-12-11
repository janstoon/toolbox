package basex_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks/wire/basex"
)

func TestEncode(t *testing.T) {
	text := "The universe we observe has precisely the properties we should expect if there is, " +
		"at bottom, no design, no purpose, no evil, no good, nothing but blind, pitiless indifference."

	b64BuiltIn := base64.StdEncoding
	b64 := basex.NewEndec("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	padChar := '%'
	b64WithPadChar := b64.WithPadding(padChar)
	b64BuiltInWithPadChar := b64BuiltIn.WithPadding(padChar)

	assert.Equal(t, []byte(b64BuiltIn.EncodeToString([]byte("Pouyan"))), b64.Encode([]byte("Pouyan")))
	for i := 0; i < len(text); i++ {
		input := []byte(text[:i+1])
		assert.Equal(t, []byte(b64BuiltIn.EncodeToString(input)), b64.Encode(input))
		assert.Equal(t, []byte(b64BuiltInWithPadChar.EncodeToString(input)), b64WithPadChar.Encode(input))
	}
}
