package bricks_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
	_ "github.com/janstoon/toolbox/bricks/countries/iran"
)

func ExampleParseNationalIdentityNumber() {
	country := bricks.LookupCountryByIsoAlphaTwoCode("IR")

	nin, err := bricks.ParseNationalIdentityNumber("۰۴۹-۹۳۷ ۰۸۹۹",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(nin.String())
	fmt.Println(nin.Country.EnglishName.Short)
	fmt.Println(bricks.MustParseNationalIdentityNumber("0499370899",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	))

	_, err = bricks.ParseNationalIdentityNumber("0499370899")
	fmt.Println(errors.Is(err, bricks.ErrNationalIdentityNumberUnknownCountry))

	// Output:
	// 0499370899
	// Iran
	// 0499370899
	// true
}

func TestParseNationalIdentityNumberArabicIndicDigits(t *testing.T) {
	country := bricks.LookupCountryByIsoAlphaTwoCode("IR")
	nin, err := bricks.ParseNationalIdentityNumber("٠٤٩٩٣٧٠٨٩٩",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	)
	require.NoError(t, err)
	assert.Equal(t, "0499370899", nin.String())
}

func TestParseNationalIdentityNumberUnknownCountry(t *testing.T) {
	nin, err := bricks.ParseNationalIdentityNumber("0499370899")
	require.ErrorIs(t, err, bricks.ErrInvalidArgument)
	require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberUnknownCountry)
	assert.Nil(t, nin)
	require.Panics(t, func() {
		bricks.MustParseNationalIdentityNumber("0499370899")
	})

	nin, err = bricks.ParseNationalIdentityNumber("0499370899",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(nil),
	)
	require.ErrorIs(t, err, bricks.ErrInvalidArgument)
	require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberUnknownCountry)
	assert.Nil(t, nin)

	nin, err = bricks.ParseNationalIdentityNumber("0499370899",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(&neverland),
	)
	require.ErrorIs(t, err, bricks.ErrInvalidArgument)
	require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberUnknownCountry)
	assert.Nil(t, nin)
	require.Panics(t, func() {
		bricks.MustParseNationalIdentityNumber("0499370899",
			bricks.ParseNationalIdentityNumberWithDefaultCountry(&neverland),
		)
	})
}
