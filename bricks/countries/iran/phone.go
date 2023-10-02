package iran

import (
	"errors"

	"github.com/janstoon/toolbox/bricks"
)

func phoneNumberLocator(localNumber string) (*bricks.NetworkOperator, error) {
	op := operators.BestMatch(localNumber)
	if op == nil {
		return nil, errors.Join(bricks.ErrInvalidInput, bricks.ErrUnknownNetworkOperator)
	}

	return op, nil
}

func init() {
	setupOperators()
	bricks.RegisterPhoneNumberLocator(iran.Codes.Telephone, phoneNumberLocator)
}

var (
	operators = bricks.Trie[string, rune, bricks.NetworkOperator](func(s string) []rune {
		return []rune(s)
	})

	oprTci = bricks.NetworkOperator{
		Name:    "TCI",
		Mobile:  false,
		Virtual: false,
	}
	oprMci = bricks.NetworkOperator{
		Name:    "MCI",
		Mobile:  true,
		Virtual: false,
	}
	oprMtnIrancellMobile = bricks.NetworkOperator{
		Name:    "MTN",
		Mobile:  true,
		Virtual: false,
	}
	oprMtnIrancellTdLte = bricks.NetworkOperator{
		Name:    "MTN",
		Mobile:  false,
		Virtual: false,
	}
	oprRightel = bricks.NetworkOperator{
		Name:    "Rightel",
		Mobile:  true,
		Virtual: false,
	}
	oprTaliya = bricks.NetworkOperator{
		Name:    "Talia",
		Mobile:  true,
		Virtual: false,
	}
	oprSpadan = bricks.NetworkOperator{
		Name:    "Spadan",
		Mobile:  true,
		Virtual: false,
	}
	oprTeleKish = bricks.NetworkOperator{
		Name:    "Kish",
		Mobile:  true,
		Virtual: false,
	}
	oprShatel = bricks.NetworkOperator{
		Name:    "Shatel",
		Mobile:  true,
		Virtual: true,
	}
	oprApTel = bricks.NetworkOperator{
		Name:    "ApTel",
		Mobile:  true,
		Virtual: true,
	}
	oprSamanTel = bricks.NetworkOperator{
		Name:    "SamanTel",
		Mobile:  true,
		Virtual: true,
	}
	oprLotusTel = bricks.NetworkOperator{
		Name:    "LotusTel",
		Mobile:  true,
		Virtual: true,
	}
	oprArianTel = bricks.NetworkOperator{
		Name:    "ArianTel",
		Mobile:  true,
		Virtual: true,
	}
)

func setupOperators() {
	operators.Put("21", oprTci)

	operators.Put("911", oprMci)
	operators.Put("912", oprMci)
	operators.Put("913", oprMci)
	operators.Put("914", oprMci)
	operators.Put("915", oprMci)
	operators.Put("916", oprMci)
	operators.Put("917", oprMci)
	operators.Put("918", oprMci)
	operators.Put("919", oprMci)

	operators.Put("990", oprMci)
	operators.Put("991", oprMci)
	operators.Put("992", oprMci)
	operators.Put("993", oprMci)
	operators.Put("994", oprMci)
	operators.Put("995", oprMci)
	operators.Put("996", oprMci)

	operators.Put("930", oprMtnIrancellMobile)
	operators.Put("933", oprMtnIrancellMobile)
	operators.Put("935", oprMtnIrancellMobile)
	operators.Put("936", oprMtnIrancellMobile)
	operators.Put("937", oprMtnIrancellMobile)
	operators.Put("938", oprMtnIrancellMobile)
	operators.Put("939", oprMtnIrancellMobile)
	operators.Put("901", oprMtnIrancellMobile)
	operators.Put("902", oprMtnIrancellMobile)
	operators.Put("903", oprMtnIrancellMobile)
	operators.Put("904", oprMtnIrancellMobile)
	operators.Put("905", oprMtnIrancellMobile)

	operators.Put("941", oprMtnIrancellTdLte)

	operators.Put("920", oprRightel)
	operators.Put("921", oprRightel)
	operators.Put("922", oprRightel)
	operators.Put("923", oprRightel)

	operators.Put("932", oprTaliya)

	operators.Put("931", oprSpadan)

	operators.Put("934", oprTeleKish)

	operators.Put("99810", oprShatel)
	operators.Put("99811", oprShatel)
	operators.Put("99812", oprShatel)
	operators.Put("99813", oprShatel)
	operators.Put("99814", oprShatel)
	operators.Put("99815", oprShatel)
	operators.Put("99816", oprShatel)
	operators.Put("99817", oprShatel)

	operators.Put("99910", oprApTel)
	operators.Put("99911", oprApTel)
	operators.Put("99913", oprApTel)
	operators.Put("99914", oprApTel)

	operators.Put("99999", oprSamanTel)
	operators.Put("9999", oprSamanTel)

	operators.Put("9990", oprLotusTel)

	operators.Put("9998", oprArianTel)
}
