package bricks

import (
	"errors"
	"fmt"
	"regexp"
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
	Mobile          bool
	Prepaid         bool
	DefaultOperator NetworkOperator
}

func (pn PhoneNumber) String() string {
	return fmt.Sprintf("+%s", pn.full)
}

type NetworkOperator struct {
	Name    string
	Virtual bool
}

type PhoneNumberMetadata struct {
	Mobile   bool
	Prepaid  bool
	Operator NetworkOperator
}

type (
	PhoneNumberResolver func(localNumber string) (*PhoneNumberMetadata, error)

	phoneNumberResolverCatalogue struct {
		countryTelCode string
		resolver       PhoneNumberResolver
	}
)

type phoneNumberResolverBank struct {
	sync.RWMutex

	byCountryTelCode *TrieNode[string, rune, phoneNumberResolverCatalogue]
}

func (bank *phoneNumberResolverBank) push(c phoneNumberResolverCatalogue) {
	bank.Lock()
	defer bank.Unlock()

	bank.byCountryTelCode.Put(c.countryTelCode, c)
}

func (bank *phoneNumberResolverBank) matchByCountryTelCode(number string) *phoneNumberResolverCatalogue {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byCountryTelCode.BestMatch(number)
}

var phoneNumberResolvers = phoneNumberResolverBank{
	byCountryTelCode: Trie[string, rune, phoneNumberResolverCatalogue](func(s string) []rune {
		return []rune(s)
	}),
}

func ParsePhoneNumber(number string) (*PhoneNumber, error) {
	number = sanitizePhoneNumber(number)

	var catalogue phoneNumberResolverCatalogue
	if c := phoneNumberResolvers.matchByCountryTelCode(number); c == nil {
		return nil, errors.Join(ErrInvalidInput, ErrPhoneNumberUnknownCountry)
	} else {
		catalogue = tricks.PtrVal(c)
	}

	localNumber := strings.TrimPrefix(strings.TrimPrefix(number, catalogue.countryTelCode), "0")
	meta, err := catalogue.resolver(localNumber)
	if err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	return &PhoneNumber{
		full:            fmt.Sprintf("%s%s", catalogue.countryTelCode, localNumber),
		Country:         LookupCountryByTelephoneCode(catalogue.countryTelCode),
		Mobile:          meta.Mobile,
		DefaultOperator: meta.Operator,
	}, nil
}

func MustParsePhoneNumber(number string) PhoneNumber {
	pn, err := ParsePhoneNumber(number)
	if err != nil {
		panic(err)
	}

	return tricks.PtrVal(pn)
}

var ptnNonDigits = regexp.MustCompile(`\D`)

func sanitizePhoneNumber(src string) string {
	return strings.TrimPrefix(ptnNonDigits.ReplaceAllString(src, ""), "00")
}

func RegisterPhoneNumberResolver(countryTelCode string, resolver PhoneNumberResolver) {
	phoneNumberResolvers.push(phoneNumberResolverCatalogue{
		countryTelCode: countryTelCode,
		resolver:       resolver,
	})
}
