package basex

import (
	"fmt"
	"math"

	"github.com/janstoon/toolbox/tricks/mathx"
)

const (
	StdPadding rune = '='
	NoPadding  rune = -1
)

type Endec struct {
	base    int
	charset []byte

	padChar rune
	strict  bool

	srcBlockBits int
	dstBlockBits int
	mask         byte
}

func NewEndec(charset string) *Endec {
	for _, char := range charset {
		if char == '\n' || char == '\r' {
			panic(fmt.Errorf("charset contains illegal character: %q", char))
		}
	}

	enc := &Endec{
		base:    len(charset),
		padChar: StdPadding,

		srcBlockBits: 8,
	}
	enc.charset = make([]byte, enc.base)
	copy(enc.charset, charset)
	enc.dstBlockBits = int(math.Ceil(math.Log2(float64(enc.base))))
	enc.mask = byte((math.Pow(2, float64(enc.dstBlockBits))) - 1)

	return enc
}

func (enc Endec) WithPadding(padding rune) *Endec {
	if padding == '\r' || padding == '\n' || padding > 0xff {
		panic(fmt.Errorf("invalid padding: %q", padding))
	}

	for i := 0; i < len(enc.charset); i++ {
		if rune(enc.charset[i]) == padding {
			panic(fmt.Errorf("padding contained in alphabet: %q", padding))
		}
	}

	enc.padChar = padding

	return &enc
}

func (enc Endec) WithoutPadding() *Endec {
	return enc.WithPadding(NoPadding)
}

func (enc Endec) Strict() *Endec {
	enc.strict = true

	return &enc
}

func (enc Endec) Encode(src []byte) []byte {
	pads := enc.pads(src)

	dl := enc.encLength(src)
	dst := make([]byte, dl)
	if dl*enc.dstBlockBits%enc.srcBlockBits > 0 {
		src = append(src, 0)
	}

	for di := range dst {
		var xtet byte
		si := (di * enc.dstBlockBits) / enc.srcBlockBits

		shift := (si+1)*enc.srcBlockBits - (di+1)*enc.dstBlockBits
		if shift < 0 {
			shift = -shift
			xtet = (src[si] & (enc.mask >> shift)) << shift
			shift = enc.srcBlockBits - shift
			si++
		}

		xtet |= (src[si] & (enc.mask << shift)) >> shift
		dst[di] = enc.charset[xtet]
	}

	return append(dst, pads...)
}

func (enc Endec) Decode(src []byte) ([]byte, error) {
	panic("not implemented")
}

func (enc Endec) encLength(src []byte) int {
	return int(math.Ceil(float64(len(src)*enc.srcBlockBits) / float64(enc.dstBlockBits)))
}

func (enc Endec) pads(src []byte) []byte {
	if enc.padChar == NoPadding {
		return nil
	}
	perfectBlockLength := mathx.Lcm(enc.srcBlockBits, enc.dstBlockBits)
	n := (perfectBlockLength - (len(src)*enc.srcBlockBits)%perfectBlockLength) % perfectBlockLength / enc.srcBlockBits

	pads := make([]byte, n)
	padChar := byte(enc.padChar)
	for k := range pads {
		pads[k] = padChar
	}

	return pads
}

func (enc Endec) decLength(src []byte) int {
	return int(math.Floor(float64(len(src)*enc.dstBlockBits) / float64(enc.srcBlockBits)))
}
