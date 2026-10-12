package bricks_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/janstoon/toolbox/bricks"
)

func ExampleParseInternationalBankAccountNumber() {
	bricks.RegisterCountry(bricks.Country{
		EnglishName: bricks.CountryName{Short: "Exampleland"},
		Codes:       bricks.CountryCode{IsoAlphaTwo: "ZZ"},
	})
	bricks.RegisterBbanValidator("ZZ", func(bban string) error {
		if len(bban) != 22 {
			return bricks.ErrIbanInvalidBban
		}

		return nil
	})

	iban, err := bricks.ParseInternationalBankAccountNumber("ZZ440143022491272900023641")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(iban)
	fmt.Println(iban.Country.EnglishName.Short)

	// Output:
	// ZZ440143022491272900023641
	// Exampleland
}

func ExampleParsePhoneNumber() {
	bricks.RegisterPhoneNumberResolver("999", func(localNumber string) (*bricks.PhoneNumberMetadata, error) {
		return &bricks.PhoneNumberMetadata{
			Mobile: strings.HasPrefix(localNumber, "123"),
			Operator: bricks.NetworkOperator{
				Name: "NeverTel",
			},
		}, nil
	})

	pn, err := bricks.ParsePhoneNumber("+999 123-456-789")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(pn)
	fmt.Println(pn.Mobile, pn.DefaultOperator.Name)

	// Output:
	// +999123456789
	// true NeverTel
}

func ExampleTimespan_Contains() {
	start := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	span := bricks.Timespan{Start: start, End: start.Add(time.Hour)}

	fmt.Println(span.Contains(start.Add(time.Minute)))
	fmt.Println(span.Contains(start))

	// Output:
	// true
	// false
}

func ExampleMoneyAmount() {
	amount := bricks.MoneyAmount{
		Value:    decimal.New(1250, -2),
		Currency: "EUR",
	}

	fmt.Println(amount.Currency, amount.Value.String())

	// Output:
	// EUR 12.5
}
