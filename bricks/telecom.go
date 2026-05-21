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

var EmptyPhoneNumber = PhoneNumber{}

type PhoneNumber struct {
	full            string
	Country         *Country
	Mobile          bool
	Prepaid         bool
	DefaultOperator NetworkOperator
}

func ParsePhoneNumber(number string, oo ...ParsePhoneNumberOption) (*PhoneNumber, error) {
	return tricks.ApplyOptions(new(phoneNumberParser), oo...).parse(number)
}

func MustParsePhoneNumber(number string, oo ...ParsePhoneNumberOption) PhoneNumber {
	pn, err := ParsePhoneNumber(number, oo...)
	if err != nil {
		panic(err)
	}

	return tricks.PtrVal(pn)
}

func (pn PhoneNumber) String() string {
	return fmt.Sprintf("+%s", pn.full)
}

type NetworkOperator struct {
	Name    string
	Virtual bool
}

type networkOperatorBank struct {
	sync.RWMutex

	byCountryIsoAlphaTwoCode map[string][]NetworkOperator
}

func (bank *networkOperatorBank) push(c *Country, nn ...NetworkOperator) {
	bank.Lock()
	defer bank.Unlock()

	if c == nil {
		panic(ErrNotFound)
	}

	bank.byCountryIsoAlphaTwoCode[c.Codes.IsoAlphaTwo] = append(bank.byCountryIsoAlphaTwoCode[c.Codes.IsoAlphaTwo], nn...)
}

func (bank *networkOperatorBank) lookupByCountryIsoAlphaTwoCode(code string) []NetworkOperator {
	bank.RLock()
	defer bank.RUnlock()

	country := LookupCountryByIsoAlphaTwoCode(code)
	if nil == country {
		return nil
	}

	return bank.byCountryIsoAlphaTwoCode[country.Codes.IsoAlphaTwo]
}

var networkOperators = networkOperatorBank{
	byCountryIsoAlphaTwoCode: make(map[string][]NetworkOperator),
}

func RegisterNetworkOperators(countryIsoAlphaTwoCode string, nn ...NetworkOperator) {
	networkOperators.push(LookupCountryByIsoAlphaTwoCode(countryIsoAlphaTwoCode), nn...)
}

// NetworkOperatorsByCountryCode Lists registered NetworkOperator(s) of a country where code is iso alpha-2 code
func NetworkOperatorsByCountryCode(code string) []NetworkOperator {
	return networkOperators.lookupByCountryIsoAlphaTwoCode(code)
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

var ptnNonDigits = regexp.MustCompile(`\D`)

type ParsePhoneNumberOption = tricks.Option[phoneNumberParser]
type phoneNumberParser struct {
	catalogue *phoneNumberResolverCatalogue
}

func (parser *phoneNumberParser) parse(number string) (*PhoneNumber, error) {
	number = sanitizePhoneNumber(number)

	if c := phoneNumberResolvers.matchByCountryTelCode(number); c != nil {
		parser.catalogue = c
		number = strings.TrimPrefix(number, parser.catalogue.countryTelCode)
	}

	if parser.catalogue == nil {
		return nil, errors.Join(ErrInvalidArgument, ErrPhoneNumberUnknownCountry)
	}

	localNumber := strings.TrimPrefix(number, "0")
	meta, err := parser.catalogue.resolver(localNumber)
	if err != nil {
		return nil, errors.Join(ErrInvalidArgument, err)
	}

	return &PhoneNumber{
		full:            fmt.Sprintf("%s%s", parser.catalogue.countryTelCode, localNumber),
		Country:         LookupCountryByTelephoneCode(parser.catalogue.countryTelCode),
		Mobile:          meta.Mobile,
		Prepaid:         meta.Prepaid,
		DefaultOperator: meta.Operator,
	}, nil
}

func ParsePhoneNumberWithDefaultCountry(country *Country) ParsePhoneNumberOption {
	return tricks.MutableOption[phoneNumberParser](func(s *phoneNumberParser) {
		s.catalogue = phoneNumberResolvers.matchByCountryTelCode(country.Codes.Telephone)
	})
}

func sanitizePhoneNumber(src string) string {
	return strings.TrimPrefix(ptnNonDigits.ReplaceAllString(src, ""), "00")
}

func RegisterPhoneNumberResolver(countryTelCode string, resolver PhoneNumberResolver) {
	phoneNumberResolvers.push(phoneNumberResolverCatalogue{
		countryTelCode: countryTelCode,
		resolver:       resolver,
	})
}
