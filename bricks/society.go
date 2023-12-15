package bricks

var EmptyNationalIdentityNumber = NationalIdentityNumber{}

type NationalIdentityNumber struct{}

func ParseNationalIdentityNumber(number string) (*NationalIdentityNumber, error) {
	return nil, ErrNotImplemented
}
