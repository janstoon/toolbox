package kareless_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/kareless"
	"github.com/janstoon/toolbox/kareless/std"
)

func TestSettings_UnmarshalJson(t *testing.T) {
	type (
		birth struct {
			Country string `json:"country,omitempty"`
			City    string `json:"city,omitempty"`
			Date    string `json:"date,omitempty"`
		}

		cfgFull struct {
			Name      string   `json:"name,omitempty"`
			Favorites []string `json:"favorites,omitempty"`
			MacPerson bool     `json:"macPerson,omitempty"`
			Birth     birth    `json:"birth,omitempty"`
		}
	)

	dir, err := os.MkdirTemp("", "kareless-unmarshal-json-settings-")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	fname := "tconf"
	fh, err := os.Create(path.Join(dir, fmt.Sprintf("%s.json", fname)))
	require.NoError(t, err)
	defer func() { _ = fh.Close() }()

	enc := json.NewEncoder(fh)
	err = enc.Encode(cfgFull{
		Name:      "John",
		Favorites: strings.Split("McDonald's,Pepsi,Captain America", ","),
		Birth: birth{
			Country: "USA",
			City:    "New York",
			Date:    "1986/23/10",
		},
	})
	require.NoError(t, err)

	ss := new(kareless.Settings)
	var (
		iCfgFull   cfgFull
		iName      string
		iFavorites []string
		iMacPerson bool
		iBirthCity string
		iUnknown   string
	)

	ss.Append(std.LocalEarlyLoadedSettingSource(fname, dir))

	err = ss.UnmarshalJson("", &iCfgFull)
	require.NoError(t, err)
	assert.Equal(t, cfgFull{
		Name: "John",
		Favorites: []string{
			"McDonald's",
			"Pepsi",
			"Captain America",
		},
		MacPerson: false,
		Birth: birth{
			Country: "USA",
			City:    "New York",
			Date:    "1986/23/10",
		},
	}, iCfgFull)

	err = ss.UnmarshalJson("name", &iName)
	require.NoError(t, err)
	assert.Equal(t, "John", iName)

	err = ss.UnmarshalJson("favorites", &iFavorites)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"McDonald's",
		"Pepsi",
		"Captain America",
	}, iFavorites)

	err = ss.UnmarshalJson("birth.city", &iBirthCity)
	require.NoError(t, err)
	assert.Equal(t, "New York", iBirthCity)

	err = ss.UnmarshalJson("macPerson", &iMacPerson)
	require.NoError(t, err)
	assert.False(t, iMacPerson)

	err = ss.UnmarshalJson("unknown", &iUnknown)
	require.NoError(t, err)
	assert.Equal(t, "", iUnknown)

	ss.Prepend(std.MapSettingSource{
		"name": "Pouyan",
		"favorites": []string{
			"Linkin Park",
			"The Killing Joke",
			"Inception",
		},
		"birth": map[string]string{
			"country": "Iran",
			"city":    "Tehran",
			"date":    "1991/05/11",
		},
	})

	err = ss.UnmarshalJson("", &iCfgFull)
	require.NoError(t, err)
	assert.Equal(t, cfgFull{
		Name: "Pouyan",
		Favorites: []string{
			"Linkin Park",
			"The Killing Joke",
			"Inception",
		},
		MacPerson: false,
		Birth: birth{
			Country: "Iran",
			City:    "Tehran",
			Date:    "1991/05/11",
		},
	}, iCfgFull)

	err = ss.UnmarshalJson("name", &iName)
	require.NoError(t, err)
	assert.Equal(t, "Pouyan", iName)

	err = ss.UnmarshalJson("favorites", &iFavorites)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"Linkin Park",
		"The Killing Joke",
		"Inception",
	}, iFavorites)

	err = ss.UnmarshalJson("birth.city", &iBirthCity)
	require.NoError(t, err)
	assert.Equal(t, "Tehran", iBirthCity)

	err = ss.UnmarshalJson("macPerson", &iMacPerson)
	require.NoError(t, err)
	assert.False(t, iMacPerson)

	err = ss.UnmarshalJson("unknown", &iUnknown)
	require.NoError(t, err)
	assert.Equal(t, "", iUnknown)

	ss.Append(std.MapSettingSource{
		"macPerson": true,
	})

	err = ss.UnmarshalJson("macPerson", &iMacPerson)
	require.NoError(t, err)
	assert.True(t, iMacPerson)
}
