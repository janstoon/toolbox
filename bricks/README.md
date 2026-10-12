# Bricks [![Coverage Badge][bdg-cov-bricks]][action-tests]
Includes international data types and primitive data structures known for humans and required by developers.
Full documentation and examples are available in [pkg.go.dev][go-pkg].

# Usage

```shell
go get github.com/janstoon/toolbox/bricks
```

## National identity number

A national identity number does not carry a country prefix. Pass the country with
`ParseNationalIdentityNumberWithDefaultCountry`. Country parsers are registered by ISO 3166-1 alpha-2.
A blank import of a country package installs that country's parser.

```go
import (
	"fmt"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func main() {
	iran := bricks.LookupCountryByIsoAlphaTwoCode("IR")
	nin, err := bricks.ParseNationalIdentityNumber("۰۴۹-۹۳۷ ۰۸۹۹",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(iran),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(nin) // 0499370899
}
```

Iran accepts exactly 10 digits. Spaces and dashes are removed, and Persian and Arabic-Indic digits are
mapped to Latin digits. A number whose digits are all the same, or whose check digit does not match, is rejected.
Calling `ParseNationalIdentityNumber` without a country returns `ErrNationalIdentityNumberUnknownCountry`.

## Utils
Resolvers, Parsers

# Types

## Banking
### International Bank Account Number
### Primary Account Number
### Bank Account
### Money

## Geography
### Country

## Internet
### Email

## Society
### National Identity Number

## Telecom
### Phone Number
### Network Operator

## Time
### Timespan
### Jalali

## Technical
Includes Abstract Data Types
### Graph, Tree, Heap, Treap, Trie
### Bag
### Filter
### Compensator

## Wire
Portable data
### Error
### Number

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-bricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_bricks_refs_heads_master.json
[go-pkg]: https://pkg.go.dev/github.com/janstoon/toolbox/bricks
