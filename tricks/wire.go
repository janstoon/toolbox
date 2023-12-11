package tricks

type Encoder interface {
	Encode(src []byte) []byte
}

type Decoder interface {
	Decode(src []byte) ([]byte, error)
}

type Endec interface {
	Encoder
	Decoder
}
