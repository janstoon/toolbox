package iran

import "github.com/janstoon/toolbox/bricks"

func init() {
	bricks.RegisterCountry(bricks.Country{
		EnglishName: bricks.CountryName{
			Short: "Iran",
			Full:  "The Islamic Republic of Iran",
		},
		Codes: bricks.CountryCode{
			IsoAlphaTwo:   "IR",
			IsoAlphaThree: "IRN",
			IsoNumeric:    "364",
			Telephone:     98,
		},
	})
}
