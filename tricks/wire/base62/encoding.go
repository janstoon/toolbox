package base62

import (
	"fmt"

	"github.com/janstoon/toolbox/tricks/wire/basex"
)

const (
	charsetLowers = "abcdefghijklmnopqrstuvwxyz"
	charsetUppers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigits = "0123456789"

	CharsetUpperLowerDigit = charsetUppers + charsetLowers + charsetDigits
	CharsetLowerUpperDigit = charsetLowers + charsetUppers + charsetDigits
	CharsetDigitUpperLower = charsetDigits + charsetUppers + charsetLowers
	CharsetDigitLowerUpper = charsetDigits + charsetLowers + charsetUppers
	CharsetDefault         = CharsetDigitUpperLower
)

var (
	UpperLowerDigitEndec = NewEndec(CharsetUpperLowerDigit).WithoutPadding()
	LowerUpperDigitEndec = NewEndec(CharsetLowerUpperDigit).WithoutPadding()
	DigitUpperLowerEndec = NewEndec(CharsetDigitUpperLower).WithoutPadding()
	DigitLowerUpperEndec = NewEndec(CharsetDigitLowerUpper).WithoutPadding()
	DefaultEndec         = NewEndec(CharsetDefault).WithoutPadding()
)

const (
	base = 62
)

func NewEndec(charset string) *basex.Endec {
	if len(charset) != base {
		panic(fmt.Sprintf("charset is not %d-bytes long", base))
	}

	return basex.NewEndec(charset)
}
