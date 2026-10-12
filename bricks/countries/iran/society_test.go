package iran_test

import (
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/bricks/countries/iran"
)

func TestParseNationalIdentityNumber(t *testing.T) {
	country := bricks.LookupCountryByIsoAlphaTwoCode(iran.Iran.Codes.IsoAlphaTwo)

	nin, err := bricks.ParseNationalIdentityNumber("0499370899",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	)
	require.NoError(t, err)
	require.NotNil(t, nin)
	assert.Equal(t, "0499370899", nin.String())
	assert.Equal(t, country, nin.Country)
	assert.Equal(t, tricks.PtrVal(nin), bricks.MustParseNationalIdentityNumber("0499370899",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	))

	for _, number := range []string{
		"049 937 0899",
		"049-937-0899",
		" 0499370899 ",
		"۰۴۹۹۳۷۰۸۹۹",
		"٠٤٩٩٣٧٠٨٩٩",
	} {
		nin, err = bricks.ParseNationalIdentityNumber(number,
			bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
		)
		require.NoError(t, err, number)
		require.NotNil(t, nin, number)
		assert.Equal(t, "0499370899", nin.String(), number)
		assert.Equal(t, country, nin.Country, number)
	}

	nin, err = bricks.ParseNationalIdentityNumber("0499370898",
		bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
	)
	require.ErrorIs(t, err, bricks.ErrInvalidArgument)
	require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberInvalid)
	require.ErrorContains(t, err, "check digit incorrect")
	assert.Nil(t, nin)

	for _, number := range []string{
		"0000000000",
		"1111111111",
	} {
		nin, err = bricks.ParseNationalIdentityNumber(number,
			bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
		)
		require.ErrorIs(t, err, bricks.ErrInvalidArgument, number)
		require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberInvalid, number)
		require.ErrorContains(t, err, "repeated digits", number)
		assert.Nil(t, nin, number)
	}

	for _, number := range []string{
		"049937089",
		"04993708990",
		"",
	} {
		nin, err = bricks.ParseNationalIdentityNumber(number,
			bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
		)
		require.ErrorIs(t, err, bricks.ErrInvalidArgument, number)
		require.ErrorIs(t, err, bricks.ErrNationalIdentityNumberInvalid, number)
		require.ErrorContains(t, err, "length incorrect", number)
		assert.Nil(t, nin, number)
		require.Panics(t, func() {
			bricks.MustParseNationalIdentityNumber(number,
				bricks.ParseNationalIdentityNumberWithDefaultCountry(country),
			)
		})
	}
}
