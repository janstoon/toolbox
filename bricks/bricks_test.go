package bricks_test

import "github.com/janstoon/toolbox/bricks"

var (
	neverland = bricks.Country{
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

	otherland = bricks.Country{
		EnglishName: bricks.CountryName{
			Short: "Otherland",
			Full:  "The other land",
		},
		Codes: bricks.CountryCode{
			IsoAlphaTwo:   "OT",
			IsoAlphaThree: "OTL",
			IsoNumeric:    "998",
			IocAlphaThree: "OTD",
			Telephone:     "998",
		},
	}
)

func init() {
	bricks.RegisterCountry(neverland)
	bricks.RegisterCountry(otherland)
}
