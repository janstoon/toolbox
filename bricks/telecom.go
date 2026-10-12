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
	// ErrPhoneNumberUnknownCountry is returned when no resolver matches the number
	// and the caller did not select a country.
	ErrPhoneNumberUnknownCountry = errors.New("phone number country unknown")
	// ErrUnknownNetworkOperator is returned by country resolvers that do not recognize the local number.
	ErrUnknownNetworkOperator = errors.New("network operator unknown")
)

// EmptyPhoneNumber is the zero PhoneNumber.
var EmptyPhoneNumber = PhoneNumber{}

// PhoneNumber is a telephone number in international digits, without a leading plus.
type PhoneNumber struct {
	full string
	// Country is the country registered for the calling code, or nil when that country is unknown.
	Country *Country
	// Mobile reports whether the country resolver classified the number as a mobile number.
	Mobile bool
	// Prepaid reports whether the country resolver classified the number as prepaid.
	Prepaid bool
	// DefaultOperator is the operator the country resolver assigned.
	DefaultOperator NetworkOperator
}

// ParsePhoneNumber parses number. Non-digits are removed and a leading 00 is stripped.
// The calling code is matched against resolvers registered with RegisterPhoneNumberResolver,
// unless ParsePhoneNumberWithDefaultCountry already selected one.
// On failure the error wraps ErrInvalidArgument.
func ParsePhoneNumber(number string, oo ...ParsePhoneNumberOption) (*PhoneNumber, error) {
	return tricks.ApplyOptions(new(phoneNumberParser), oo...).parse(number)
}

// MustParsePhoneNumber is like ParsePhoneNumber but panics if number cannot be parsed.
func MustParsePhoneNumber(number string, oo ...ParsePhoneNumberOption) PhoneNumber {
	pn, err := ParsePhoneNumber(number, oo...)
	if err != nil {
		panic(err)
	}

	return tricks.PtrVal(pn)
}

// String returns the number in E.164 form with a leading plus.
func (pn PhoneNumber) String() string {
	return fmt.Sprintf("+%s", pn.full)
}

// NetworkOperator is a telephone network operator in one country.
type NetworkOperator struct {
	// Name is the operator's display name.
	Name string
	// Virtual reports whether the operator is a mobile virtual network operator.
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

// RegisterNetworkOperators appends operators for a country that was already registered.
// It panics with ErrNotFound when countryIsoAlphaTwoCode is unknown.
func RegisterNetworkOperators(countryIsoAlphaTwoCode string, nn ...NetworkOperator) {
	networkOperators.push(LookupCountryByIsoAlphaTwoCode(countryIsoAlphaTwoCode), nn...)
}

// NetworkOperatorsByCountryCode Lists registered NetworkOperator(s) of a country where code is iso alpha-2 code
func NetworkOperatorsByCountryCode(code string) []NetworkOperator {
	return networkOperators.lookupByCountryIsoAlphaTwoCode(code)
}

// PhoneNumberMetadata is what a PhoneNumberResolver reports about a local number.
type PhoneNumberMetadata struct {
	// Mobile reports whether the number is a mobile number.
	Mobile bool
	// Prepaid reports whether the number is prepaid.
	Prepaid bool
	// Operator is the operator that owns the number.
	Operator NetworkOperator
}

type (
	// PhoneNumberResolver classifies a local number, without the country calling code and without a leading zero.
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

// ParsePhoneNumberOption configures ParsePhoneNumber.
type (
	ParsePhoneNumberOption = tricks.Option[phoneNumberParser]
	phoneNumberParser      struct {
		catalogue *phoneNumberResolverCatalogue
	}
)

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

// ParsePhoneNumberWithDefaultCountry selects the resolver registered for country.Codes.Telephone
// when the number itself has no matching calling-code prefix.
func ParsePhoneNumberWithDefaultCountry(country *Country) ParsePhoneNumberOption {
	return tricks.MutableOption[phoneNumberParser](func(s *phoneNumberParser) {
		s.catalogue = phoneNumberResolvers.matchByCountryTelCode(country.Codes.Telephone)
	})
}

func sanitizePhoneNumber(src string) string {
	return strings.TrimPrefix(ptnNonDigits.ReplaceAllString(src, ""), "00")
}

// RegisterPhoneNumberResolver registers resolver for a calling code without a leading plus.
// Longer codes are preferred when several codes are prefixes of the same number.
func RegisterPhoneNumberResolver(countryTelCode string, resolver PhoneNumberResolver) {
	phoneNumberResolvers.push(phoneNumberResolverCatalogue{
		countryTelCode: countryTelCode,
		resolver:       resolver,
	})
}
