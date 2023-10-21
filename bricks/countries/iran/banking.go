package iran

import (
	"errors"
	"strconv"

	"github.com/janstoon/toolbox/bricks"
)

func bbanValidator(basicBankAccountNumber string) error {
	if len(basicBankAccountNumber) != 22 {
		return errors.Join(bricks.ErrIbanInvalidBban, errors.New("bban length incorrect"))
	}

	bankIdentifier := basicBankAccountNumber[0:3]
	// accountType := basicBankAccountNumber[3]
	// accountNumber := basicBankAccountNumber[4:]

	const biCbi = 10

	if bi, err := strconv.Atoi(bankIdentifier); err != nil {
		return errors.Join(bricks.ErrIbanInvalidBban, err)
	} else if bi < biCbi {
		return errors.Join(bricks.ErrIbanInvalidBban, errors.New("invalid bank identifier"))
	}

	return nil
}

type primaryAccountNumber struct {
	Institute   [6]byte
	Product     [2]byte
	Serial      [7]byte
	CheckDigits [1]byte
}

func init() {
	bricks.RegisterBbanValidator(iran.Codes.IsoAlphaTwo, bbanValidator)
}
