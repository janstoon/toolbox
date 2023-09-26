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
	//accountType := basicBankAccountNumber[3]
	//accountNumber := basicBankAccountNumber[4:]

	const biCbi = 10

	if bi, err := strconv.Atoi(bankIdentifier); err != nil {
		return errors.Join(bricks.ErrIbanInvalidBban, err)
	} else if bi < biCbi {
		return errors.Join(bricks.ErrIbanInvalidBban, errors.New("invalid bank identifier"))
	}

	return nil
}

func init() {
	bricks.RegisterBbanValidator("IR", bbanValidator)
}
