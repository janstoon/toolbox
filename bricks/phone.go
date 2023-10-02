package bricks

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/janstoon/toolbox/tricks"
)

var (
	ErrPhoneNumberUnknownCountry = errors.New("phone number country unknown")
	ErrUnknownNetworkOperator    = errors.New("network operator unknown")
)

type PhoneNumber struct {
	full            string
	Country         *Country
	DefaultOperator NetworkOperator
}

func (pn PhoneNumber) String() string {
	return fmt.Sprintf("+%s", pn.full)
}

type NetworkOperator struct {
	Name    string
	Mobile  bool
	Virtual bool
}

type (
	PhoneNumberLocator func(localNumber string) (*NetworkOperator, error)

	phoneNumberLocatorCatalogue struct {
		countryTelCode string
		locator        PhoneNumberLocator
	}
)

type phoneNumberLocatorBank struct {
	sync.RWMutex

	byCountryTelCode *TrieNode[string, rune, phoneNumberLocatorCatalogue]
}

func (bank *phoneNumberLocatorBank) push(c phoneNumberLocatorCatalogue) {
	bank.Lock()
	defer bank.Unlock()

	bank.byCountryTelCode.Put(c.countryTelCode, c)
}

func (bank *phoneNumberLocatorBank) matchByCountryTelCode(number string) *phoneNumberLocatorCatalogue {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byCountryTelCode.BestMatch(number)
}

var phoneNumberLocators = phoneNumberLocatorBank{
	byCountryTelCode: Trie[string, rune, phoneNumberLocatorCatalogue](func(s string) []rune {
		return []rune(s)
	}),
}

func ParsePhoneNumber(number string) (*PhoneNumber, error) {
	number = sanitizePhoneNumber(number)

	var catalogue phoneNumberLocatorCatalogue
	if c := phoneNumberLocators.matchByCountryTelCode(number); c == nil {
		return nil, errors.Join(ErrInvalidInput, ErrPhoneNumberUnknownCountry)
	} else {
		catalogue = tricks.PtrVal(c)
	}

	localNumber := strings.TrimPrefix(strings.TrimPrefix(number, catalogue.countryTelCode), "0")
	operator, err := catalogue.locator(localNumber)
	if err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	return &PhoneNumber{
		full:            fmt.Sprintf("%s%s", catalogue.countryTelCode, localNumber),
		Country:         LookupCountryByTelephoneCode(catalogue.countryTelCode),
		DefaultOperator: tricks.PtrVal(operator),
	}, nil
}

func sanitizePhoneNumber(src string) string {
	return strings.TrimPrefix(
		strings.TrimPrefix(
			strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(src),
			"+",
		),
		"00",
	)
}

func RegisterPhoneNumberLocator(countryTelCode string, locator PhoneNumberLocator) {
	phoneNumberLocators.push(phoneNumberLocatorCatalogue{
		countryTelCode: countryTelCode,
		locator:        locator,
	})
}
