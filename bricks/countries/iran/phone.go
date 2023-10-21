package iran

import (
	"errors"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
)

func phoneNumberResolver(localNumber string) (*bricks.PhoneNumberMetadata, error) {
	if len(localNumber) != 10 {
		return nil, errors.Join(bricks.ErrInvalidInput, errors.New("local number length incorrect"))
	}

	md := phoneNumbersMetadata.BestMatch(localNumber)
	if md == nil {
		return nil, errors.Join(bricks.ErrInvalidInput, bricks.ErrUnknownNetworkOperator)
	}

	return md, nil
}

func init() {
	setupPhoneNumberMetadataTree()
	bricks.RegisterPhoneNumberResolver(iran.Codes.Telephone, phoneNumberResolver)
}

var phoneNumbersMetadata = bricks.Trie[string, rune, bricks.PhoneNumberMetadata](tricks.StringToRunes)

type phoneNumberMetadata struct {
	operatorSlug string
	mobile       bool
	prepaid      bool

	// provinceCodes keeps iso 3166-2 subdivision codes.
	// useful links: https://en.wikipedia.org/wiki/ISO_3166-2:IR
	provinceCodes []string

	countyName string
}

//nolint:maintidx
func setupPhoneNumberMetadataTree() {
	operatorsBySlug := map[string]bricks.NetworkOperator{
		"tci": {
			Name:    "TCI",
			Virtual: false,
		},
		"asiatech": {
			Name:    "AsiaTech",
			Virtual: false,
		},

		"mci": {
			Name:    "MCI",
			Virtual: false,
		},
		"mtn": {
			Name:    "MTN",
			Virtual: false,
		},
		"rightel": {
			Name:    "Rightel",
			Virtual: false,
		},
		"taliya": {
			Name:    "Taliya",
			Virtual: false,
		},

		// regional
		"spadan": {
			Name:    "MTCE",
			Virtual: false,
		},
		"telekish": {
			Name:    "TeleKish",
			Virtual: false,
		},

		// virtual
		"shatel": {
			Name:    "Shatel",
			Virtual: true,
		},
		"aptel": {
			Name:    "ApTel",
			Virtual: true,
		},
		"samantel": {
			Name:    "SamanTel",
			Virtual: true,
		},
		"lotustel": {
			Name:    "LotusTel",
			Virtual: true,
		},
		"ariantel": {
			Name:    "ArianTel",
			Virtual: true,
		},
	}

	metadataByPrefix := map[string]phoneNumberMetadata{
		"11": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"02"},
			countyName:    "",
		},
		"13": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"01"},
			countyName:    "",
		},
		"17": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"27"},
			countyName:    "",
		},
		"21": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"23"},
			countyName:    "",
		},
		"23": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"20"},
			countyName:    "",
		},
		"24": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"19"},
			countyName:    "",
		},
		"25": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"25"},
			countyName:    "",
		},
		"26": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"30"},
			countyName:    "",
		},
		"28": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"26"},
			countyName:    "",
		},
		"31": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"10"},
			countyName:    "",
		},
		"34": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"08"},
			countyName:    "",
		},
		"35": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"21"},
			countyName:    "",
		},
		"38": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"14"},
			countyName:    "",
		},
		"41": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"03"},
			countyName:    "",
		},
		"44": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"04"},
			countyName:    "",
		},
		"45": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"24"},
			countyName:    "",
		},
		"51": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"09"},
			countyName:    "",
		},
		"54": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"11"},
			countyName:    "",
		},
		"56": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"29"},
			countyName:    "",
		},
		"58": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"28"},
			countyName:    "",
		},
		"61": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"06"},
			countyName:    "",
		},
		"66": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"15"},
			countyName:    "",
		},
		"71": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"07"},
			countyName:    "",
		},
		"74": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"17"},
			countyName:    "",
		},
		"76": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"22"},
			countyName:    "",
		},
		"77": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"18"},
			countyName:    "",
		},
		"81": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"13"},
			countyName:    "",
		},
		"83": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"05"},
			countyName:    "",
		},
		"84": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"16"},
			countyName:    "",
		},
		"86": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"00"},
			countyName:    "",
		},
		"87": {
			operatorSlug:  "tci",
			mobile:        false,
			prepaid:       false,
			provinceCodes: []string{"12"},
			countyName:    "",
		},

		"900": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"901": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"902": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"903": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"904": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"905": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"910": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"911": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"02"},
			countyName:    "",
		},
		"912": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"23"},
			countyName:    "",
		},
		"913": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"10"},
			countyName:    "",
		},
		"914": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"03"},
			countyName:    "",
		},
		"915": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"09"},
			countyName:    "",
		},
		"916": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"06"},
			countyName:    "",
		},
		"917": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"07"},
			countyName:    "",
		},
		"918": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"13"},
			countyName:    "",
		},
		"919": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: []string{"23"},
			countyName:    "",
		},

		"920": {
			operatorSlug:  "rightel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"921": {
			operatorSlug:  "rightel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"922": {
			operatorSlug:  "rightel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"923": {
			operatorSlug:  "rightel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"930": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"931": {
			operatorSlug:  "spadan",
			mobile:        true,
			prepaid:       true,
			provinceCodes: []string{"10"},
			countyName:    "Isfahan",
		},

		"932": {
			operatorSlug:  "taliya",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},

		"933": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"934": {
			operatorSlug:  "telekish",
			mobile:        true,
			prepaid:       false,
			provinceCodes: []string{"22"},
			countyName:    "Kish",
		},

		"935": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"936": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"937": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"938": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"939": {
			operatorSlug:  "mtn",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"941": {
			operatorSlug:  "mtn",
			mobile:        false, // fixme: it's td-lte but it's rechargeable
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"990": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"991": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"992": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"993": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"994": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"995": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"996": {
			operatorSlug:  "mci",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},

		"9422": {
			operatorSlug:  "asiatech",
			mobile:        false,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"99810": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99811": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99812": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99813": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99814": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99815": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99816": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},
		"99817": {
			operatorSlug:  "shatel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},

		"99910": {
			operatorSlug:  "aptel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"99911": {
			operatorSlug:  "aptel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"99913": {
			operatorSlug:  "aptel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"99914": {
			operatorSlug:  "aptel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},

		"9990": {
			operatorSlug:  "lotustel",
			mobile:        true,
			prepaid:       true,
			provinceCodes: nil,
			countyName:    "",
		},

		"9998": {
			operatorSlug:  "ariantel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},

		"9999": {
			operatorSlug:  "samantel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
		"99999": {
			operatorSlug:  "samantel",
			mobile:        true,
			prepaid:       false,
			provinceCodes: nil,
			countyName:    "",
		},
	}

	for prefix, meta := range metadataByPrefix {
		op, ok := operatorsBySlug[meta.operatorSlug]
		if !ok {
			panic(bricks.ErrUnknownNetworkOperator)
		}

		phoneNumbersMetadata.Put(prefix, bricks.PhoneNumberMetadata{
			Mobile:   meta.mobile,
			Prepaid:  meta.prepaid,
			Operator: op,
		})
	}
}
