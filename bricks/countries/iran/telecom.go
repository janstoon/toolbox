package iran

import (
	_ "embed"
	"errors"

	"github.com/janstoon/toolbox/tricks"
	"gopkg.in/yaml.v3"

	"github.com/janstoon/toolbox/bricks"
)

//go:embed telecom.yml
var ymlTelecomPolicy string

func init() {
	registerTelecom([]byte(ymlTelecomPolicy))
}

func registerTelecom(ymlPolicy []byte) {
	var policy struct {
		Operators map[string]struct {
			Name        string `yaml:"name"`
			Virtual     bool   `yaml:"virtual"`
			Subscribers []struct {
				Prefixes  []string `yaml:"prefixes"`
				Mobile    bool     `yaml:"mobile"`
				Prepaid   bool     `yaml:"prepaid"`
				Provinces []string `yaml:"provinces"`
				Counties  []string `yaml:"counties"`
			} `yaml:"subscribers"`
		} `yaml:"operators"`
	}

	err := yaml.Unmarshal(ymlPolicy, &policy)
	if err != nil {
		panic(err)
	}

	p2m := bricks.Trie[string, rune, bricks.PhoneNumberMetadata](tricks.StringToRunes)
	operators := make([]bricks.NetworkOperator, 0)
	for _, op := range policy.Operators {
		no := bricks.NetworkOperator{
			Name:    op.Name,
			Virtual: op.Virtual,
		}
		operators = append(operators, no)

		for _, subs := range op.Subscribers {
			meta := bricks.PhoneNumberMetadata{
				Mobile:   subs.Mobile,
				Prepaid:  subs.Prepaid,
				Operator: no,
			}

			for _, prefix := range subs.Prefixes {
				p2m.Put(prefix, meta)
			}
		}
	}

	bricks.RegisterNetworkOperators(Iran.Codes.IsoAlphaTwo, operators...)
	bricks.RegisterPhoneNumberResolver(Iran.Codes.Telephone,
		func(localNumber string) (*bricks.PhoneNumberMetadata, error) {
			if len(localNumber) != 10 {
				return nil, errors.Join(bricks.ErrInvalidInput, errors.New("local number length incorrect"))
			}

			meta := p2m.BestMatch(localNumber)
			if meta == nil {
				return nil, errors.Join(bricks.ErrInvalidInput, bricks.ErrUnknownNetworkOperator)
			}

			return meta, nil
		},
	)
}
