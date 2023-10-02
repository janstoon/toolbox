package bricks_test

import "github.com/janstoon/toolbox/bricks"

var neverland = bricks.Country{
	EnglishName: bricks.CountryName{
		Short: "Neverland",
		Full:  "The Promised Neverland",
	},
	Codes: bricks.CountryCode{
		IsoAlphaTwo:   "NV",
		IsoAlphaThree: "NVL",
		IsoNumeric:    "999",
		IocAlphaThree: "NVD",
		Telephone:     "999",
	},
}

func init() {
	bricks.RegisterCountry(neverland)
}
