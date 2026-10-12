package iran

import (
	"errors"

	"github.com/janstoon/toolbox/bricks"
)

func init() {
	bricks.RegisterNationalIdentityNumberParser(Iran.Codes.IsoAlphaTwo, parseNationalIdentityNumber)
}

func parseNationalIdentityNumber(number string) error {
	const length = 10
	if len(number) != length {
		return errors.Join(bricks.ErrNationalIdentityNumberInvalid, errors.New("length incorrect"))
	}

	digits := make([]int, length)
	for i := 0; i < length; i++ {
		if number[i] < '0' || number[i] > '9' {
			return errors.Join(bricks.ErrNationalIdentityNumberInvalid, errors.New("non-digit"))
		}

		digits[i] = int(number[i] - '0')
	}

	repeated := true
	for i := 1; i < length; i++ {
		if digits[i] != digits[0] {
			repeated = false

			break
		}
	}

	if repeated {
		return errors.Join(bricks.ErrNationalIdentityNumberInvalid, errors.New("repeated digits"))
	}

	sum := 0
	for i := 0; i < length-1; i++ {
		sum += digits[i] * (length - i)
	}

	remainder := sum % 11
	check := remainder
	if remainder >= 2 {
		check = 11 - remainder
	}

	if digits[length-1] != check {
		return errors.Join(bricks.ErrNationalIdentityNumberInvalid, errors.New("check digit incorrect"))
	}

	return nil
}
