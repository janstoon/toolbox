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

Country packages register parsers and validators from `init`. Callers look values up with
`LookupCountryByIsoAlphaTwoCode` and `LookupCountryByTelephoneCode`.

# Types

## Banking

`ParseInternationalBankAccountNumber` checks length, the mod-97 checksum, and a BBAN validator
registered for the country code. `ParsePrimaryAccountNumber` accepts 8 to 19 characters.
`MoneyAmount` pairs a `decimal.Decimal` with a currency code.

## Geography

`Country` stores English names and ISO, IOC, and telephone codes. `RegisterCountry` publishes one.

## Internet

`Email` is a placeholder for an email address.

## Society

`NationalIdentityNumber` is documented above.

## Telecom

`ParsePhoneNumber` strips punctuation and a leading `00`, then asks the resolver registered for
the calling code to classify the local number. `RegisterNetworkOperators` lists operators for a
country that is already registered.

## Time

`Timespan.Contains` reports whether an instant is strictly inside the span.

## Technical

`Graph` is an adjacency map with breadth-first search. `Trie` stores values by a caller-supplied
key splitter and can return the longest stored prefix. `Queue` is a goroutine-safe FIFO.
`Registry` is a name-to-value map that rejects duplicates. `Bag` pulls items until nothing remains.
`Compensator` undoes a failed step; `CompensatorAsError` carries one inside an error.

## Wire

Coded errors use the same numbering as gRPC status codes. `ErrorWithCode` attaches a code.
`ParseError` classifies cancellations and timeouts. `MessageEnvelope` carries an id, a retry count,
and a body.

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-bricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_bricks_refs_heads_master.json
[go-pkg]: https://pkg.go.dev/github.com/janstoon/toolbox/bricks
