package iran

import "github.com/janstoon/toolbox/bricks"

var iran = bricks.Country{
	EnglishName: bricks.CountryName{
		Short: "Iran",
		Full:  "The Islamic Republic of Iran",
	},
	Codes: bricks.CountryCode{
		IsoAlphaTwo:   "IR",
		IsoAlphaThree: "IRN",
		IsoNumeric:    "364",
		IocAlphaThree: "IRI",
		Telephone:     98,
	},
}

func init() {
	bricks.RegisterCountry(iran)
}
