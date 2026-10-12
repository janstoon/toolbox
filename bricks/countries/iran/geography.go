// Package iran registers Iran (ISO 3166-1 alpha-2 IR) with bricks.
// Importing the package registers the country, its IBAN BBAN validator, phone-number resolver,
// and national identity number parser.
package iran

import "github.com/janstoon/toolbox/bricks"

var Iran = bricks.Country{
	EnglishName: bricks.CountryName{
		Short: "Iran",
		Full:  "The Islamic Republic of Iran",
	},
	Codes: bricks.CountryCode{
		IsoAlphaTwo:   "IR",
		IsoAlphaThree: "IRN",
		IsoNumeric:    "364",
		IocAlphaThree: "IRI",
		Telephone:     "98",
	},
}

func init() {
	bricks.RegisterCountry(Iran)
}
