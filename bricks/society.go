package bricks

import (
	"errors"
	"strings"
	"sync"

	"github.com/janstoon/toolbox/tricks"
)

var (
	// ErrNationalIdentityNumberUnknownCountry is returned when no country parser can run.
	// ParseNationalIdentityNumber returns it, joined with ErrInvalidArgument, when the caller
	// passes no country or passes a country that has no parser registered for its ISO 3166-1
	// alpha-2 code. The country is never inferred from the digits.
	ErrNationalIdentityNumberUnknownCountry = errors.New("national identity number country unknown")

	// ErrNationalIdentityNumberInvalid is returned when a registered country parser rejects
	// the sanitized number. ParseNationalIdentityNumber joins it with ErrInvalidArgument.
	ErrNationalIdentityNumberInvalid = errors.New("national identity number invalid")
)

// EmptyNationalIdentityNumber is the zero NationalIdentityNumber.
var EmptyNationalIdentityNumber = NationalIdentityNumber{}

// NationalIdentityNumber is a national identity number in canonical digits.
// Country is the country the caller selected. It is not inferred from the number.
type NationalIdentityNumber struct {
	full string

	Country *Country
}

// String returns the canonical digits.
func (nin NationalIdentityNumber) String() string {
	return nin.full
}

// ParseNationalIdentityNumber parses number using the country selected by oo.
// Spaces and dashes are removed, and Persian digits (۰–۹) and Arabic-Indic digits (٠–٩)
// are mapped to Latin digits before the country parser runs.
//
// With no country option, or when that country has no parser registered by
// RegisterNationalIdentityNumberParser, it returns an error that wraps ErrInvalidArgument
// and ErrNationalIdentityNumberUnknownCountry.
func ParseNationalIdentityNumber(
	number string, oo ...ParseNationalIdentityNumberOption,
) (*NationalIdentityNumber, error) {
	return tricks.ApplyOptions(new(nationalIdentityNumberParser), oo...).parse(number)
}

// MustParseNationalIdentityNumber is like ParseNationalIdentityNumber but panics if number cannot be parsed.
func MustParseNationalIdentityNumber(number string, oo ...ParseNationalIdentityNumberOption) NationalIdentityNumber {
	nin, err := ParseNationalIdentityNumber(number, oo...)
	if err != nil {
		panic(err)
	}

	return tricks.PtrVal(nin)
}

// ParseNationalIdentityNumberOption configures ParseNationalIdentityNumber.
type ParseNationalIdentityNumberOption = tricks.Option[nationalIdentityNumberParser]

type nationalIdentityNumberParser struct {
	country *Country
}

func (parser *nationalIdentityNumberParser) parse(number string) (*NationalIdentityNumber, error) {
	number = sanitizeNationalIdentityNumber(number)

	if parser.country == nil {
		return nil, errors.Join(ErrInvalidArgument, ErrNationalIdentityNumberUnknownCountry)
	}

	validate := nationalIdentityNumberParsers.lookupByCountryCode(parser.country.Codes.IsoAlphaTwo)
	if validate == nil {
		return nil, errors.Join(ErrInvalidArgument, ErrNationalIdentityNumberUnknownCountry)
	}

	if err := validate(number); err != nil {
		return nil, errors.Join(ErrInvalidArgument, err)
	}

	return &NationalIdentityNumber{
		full:    number,
		Country: parser.country,
	}, nil
}

// ParseNationalIdentityNumberWithDefaultCountry selects the country whose parser validates number.
// A national identity number does not carry a country prefix, so the caller must pass one.
func ParseNationalIdentityNumberWithDefaultCountry(country *Country) ParseNationalIdentityNumberOption {
	return tricks.MutableOption[nationalIdentityNumberParser](func(parser *nationalIdentityNumberParser) {
		parser.country = country
	})
}

func sanitizeNationalIdentityNumber(src string) string {
	var b strings.Builder
	b.Grow(len(src))

	for _, r := range src {
		switch {
		case r == ' ' || r == '-':
			continue
		case r >= '۰' && r <= '۹':
			b.WriteByte(byte('0' + r - '۰'))
		case r >= '٠' && r <= '٩':
			b.WriteByte(byte('0' + r - '٠'))
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}

type (
	// NationalIdentityNumberParser validates a sanitized national identity number for one country.
	// number has had spaces and dashes removed, and Persian and Arabic-Indic digits mapped to Latin digits.
	NationalIdentityNumberParser func(number string) error

	nationalIdentityNumberParserCatalogue struct {
		countryCode string
		parser      NationalIdentityNumberParser
	}
)

type nationalIdentityNumberParserBank struct {
	sync.RWMutex

	byCountryCode   map[string]NationalIdentityNumberParser
	normCountryCode func(string) string
}

func (bank *nationalIdentityNumberParserBank) push(c nationalIdentityNumberParserCatalogue) {
	bank.Lock()
	defer bank.Unlock()

	tricks.InsertIfNotExist(bank.byCountryCode, bank.normCountryCode(c.countryCode), c.parser)
}

func (bank *nationalIdentityNumberParserBank) lookupByCountryCode(code string) NationalIdentityNumberParser {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byCountryCode[bank.normCountryCode(code)]
}

var nationalIdentityNumberParsers = nationalIdentityNumberParserBank{
	byCountryCode:   make(map[string]NationalIdentityNumberParser),
	normCountryCode: strings.ToUpper,
}

// RegisterNationalIdentityNumberParser registers parser for an ISO 3166-1 alpha-2 country code.
// Registering the same country twice panics. A blank import of a country package is the usual way
// to install its parser.
func RegisterNationalIdentityNumberParser(countryIsoAlphaTwoCode string, parser NationalIdentityNumberParser) {
	nationalIdentityNumberParsers.push(nationalIdentityNumberParserCatalogue{
		countryCode: countryIsoAlphaTwoCode,
		parser:      parser,
	})
}
