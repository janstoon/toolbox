package bricks

import (
	"strings"
	"sync"

	"github.com/janstoon/toolbox/tricks"
)

// Country
// useful links: http://download.geonames.org/export/dump/countryInfo.txt,
// https://github.com/mledoze/countries/blob/master/countries.json
type Country struct {
	EnglishName CountryName
	Codes       CountryCode
}

type CountryName struct {
	Short string
	Full  string
}

// CountryCode contains country unique codes in different standards
type CountryCode struct {
	// IsoAlphaTwo contains ISO 3166-1 alpha-2 code
	// useful links: https://en.wikipedia.org/wiki/ISO_3166-1 and https://www.iban.com/country-codes
	IsoAlphaTwo string

	// IsoAlphaThree contains ISO 3166-1 alpha-3 code
	// useful links: https://en.wikipedia.org/wiki/ISO_3166-1 and https://www.iban.com/country-codes
	IsoAlphaThree string

	// IsoNumeric contains ISO 3166-1 numeric code
	// useful links: https://en.wikipedia.org/wiki/ISO_3166-1 and https://www.iban.com/country-codes
	IsoNumeric string

	// Telephone is the calling code
	// useful links: https://en.wikipedia.org/wiki/List_of_country_calling_codes
	Telephone int
}

type countryBank struct {
	sync.RWMutex

	byIsoAlphaTwoCode map[string]*Country
	normIsoAlphaTwo   func(string) string
	byTelephoneCode   map[int]*Country
}

func (bank *countryBank) push(c Country) {
	bank.Lock()
	defer bank.Unlock()

	tricks.InsertIfNotExist(bank.byIsoAlphaTwoCode, bank.normIsoAlphaTwo(c.Codes.IsoAlphaTwo), &c)
	tricks.InsertIfNotExist(bank.byTelephoneCode, c.Codes.Telephone, &c)
}

func (bank *countryBank) lookupByIsoAlphaTwoCode(code string) *Country {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byIsoAlphaTwoCode[bank.normIsoAlphaTwo(code)]
}

func (bank *countryBank) lookupByTelephoneCode(code int) *Country {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byTelephoneCode[code]
}

var countries = countryBank{
	byIsoAlphaTwoCode: make(map[string]*Country),
	normIsoAlphaTwo:   strings.ToUpper,
	byTelephoneCode:   make(map[int]*Country),
}

func LookupCountryByIsoAlphaTwoCode(code string) *Country {
	return countries.lookupByIsoAlphaTwoCode(code)
}

func LookupCountryByTelephoneCode(code int) *Country {
	return countries.lookupByTelephoneCode(code)
}

func RegisterCountry(c Country) {
	countries.push(c)
}
