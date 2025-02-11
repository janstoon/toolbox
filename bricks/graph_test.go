package bricks_test

import (
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestGraph_Get(t *testing.T) {
	iran := iranCities()

	assert.Nil(t, iran["IRUNK"])
	assert.Equal(t, "Tehran", iran["Tehran"].Value)
	assert.Equal(t, "Karaj", iran["Karaj"].Value)
	assert.Equal(t, "Rasht", iran["Rasht"].Value)
}

func TestGraph_BFS(t *testing.T) {
	iran := iranCities()

	// Gorgan to Bandare Parsian
	assert.Equal(t,
		[]string{"Gorgan", "Sari", "Tehran", "Qom", "Isfahan", "Yasuj", "Shiraz", "Bandare Parsian"},
		iran.BreadthFirstSearch("Gorgan", "Bandare Parsian"),
	)
}

func iranCities() bricks.Graph[string, int, string] {
	g := make(bricks.Graph[string, int, string])
	g.Add("Tehran", "Tehran")
	g.Add("Karaj", "Karaj")
	g.Add("Qazvin", "Qazvin")
	g.Add("Qom", "Qom")
	g.Add("Zanjan", "Zanjan")
	g.Add("Rasht", "Rasht")
	g.Add("Ardabil", "Ardabil")
	g.Add("Tabriz", "Tabriz")
	g.Add("Urmia", "Urmia")
	g.Add("Gorgan", "Gorgan")
	g.Add("Sari", "Sari")
	g.Add("Mashad", "Mashad")
	g.Add("Hamadan", "Hamadan")
	g.Add("Saqqez", "Saqqez")
	g.Add("Sanandaj", "Sanandaj")
	g.Add("Kermanshah", "Kermanshah")
	g.Add("Khorramabad", "Khorramabad")
	g.Add("Kashan", "Kashan")
	g.Add("Isfahan", "Isfahan")
	g.Add("Shahr-e Kord", "Shahr-e Kord")
	g.Add("Garmsar", "Garmsar")
	g.Add("Semnan", "Semnan")
	g.Add("Damghan", "Damghan")
	g.Add("Shahrud", "Shahrud")
	g.Add("Sabzevar", "Sabzevar")
	g.Add("Neyshabur", "Neyshabur")
	g.Add("Ardestan", "Ardestan")
	g.Add("Naeen", "Naeen")
	g.Add("Yazd", "Yazd")
	g.Add("Arak", "Arak")
	g.Add("Yasuj", "Yasuj")
	g.Add("Kazerun", "Kazerun")
	g.Add("Shiraz", "Shiraz")
	g.Add("Abadan", "Abadan")
	g.Add("Ahwaz", "Ahwaz")
	g.Add("Dezful", "Dezful")
	g.Add("Kerman", "Kerman")
	g.Add("Sirjan", "Sirjan")
	g.Add("Bushehr", "Bushehr")
	g.Add("Bandar Abbas", "Bandar Abbas")
	g.Add("Bandare Parsian", "Bandare Parsian")

	g.Connect("Tehran", "Karaj", "Qom", "Garmsar", "Sari")

	// North
	g.Connect("Rasht", "Ardabil", "Sari")
	g.Connect("Sari", "Semnan", "Gorgan")
	g.Connect("Gorgan", "Damghan")

	// West
	g.Connect("Karaj", "Qazvin")
	g.Connect("Qazvin", "Zanjan", "Rasht")
	g.Connect("Zanjan", "Tabriz", "Sanandaj")
	g.Connect("Tabriz", "Ardabil", "Urmia", "Saqqez")
	g.Connect("Urmia", "Saqqez")
	g.Connect("Saqqez", "Sanandaj")
	g.Connect("Sanandaj", "Kermanshah")
	g.Connect("Hamadan", "Sanandaj", "Kermanshah")
	g.Connect("Kermanshah", "Khorramabad")
	g.Connect("Khorramabad", "Dezful")
	g.Connect("Dezful", "Ahwaz")
	g.Connect("Ahwaz", "Abadan")
	g.Connect("Abadan", "Bushehr")

	// East
	g.Connect("Garmsar", "Semnan")
	g.Connect("Semnan", "Damghan")
	g.Connect("Damghan", "Shahrud")
	g.Connect("Shahrud", "Sabzevar")
	g.Connect("Sabzevar", "Mashad")

	// South
	g.Connect("Qom", "Arak", "Kashan", "Isfahan")
	g.Connect("Arak", "Hamadan", "Khorramabad")
	g.Connect("Kashan", "Isfahan", "Ardestan")
	g.Connect("Ardestan", "Naeen")
	g.Connect("Naeen", "Isfahan", "Yazd")
	g.Connect("Yazd", "Kerman", "Sirjan")
	g.Connect("Sirjan", "Bandar Abbas")
	g.Connect("Isfahan", "Shahr-e Kord", "Yasuj")
	g.Connect("Shahr-e Kord", "Yasuj", "Khorramabad")
	g.Connect("Yasuj", "Kazerun", "Shiraz")
	g.Connect("Kazerun", "Shiraz", "Bushehr")
	g.Connect("Shiraz", "Bushehr", "Bandare Parsian", "Bandar Abbas", "Sirjan")
	g.Connect("Kerman", "Sirjan")
	g.Connect("Bushehr", "Bandare Parsian")
	g.Connect("Bandar Abbas", "Bandare Parsian")

	return g
}

func TestTrie_Get(t *testing.T) {
	trie := bricks.Trie[string, rune, string](tricks.StringToRunes)

	assert.Nil(t, trie.Get("1"))
	assert.Nil(t, trie.Get("18"))
	assert.Nil(t, trie.Get("1800"))
	assert.Nil(t, trie.Get("181"))
	assert.Nil(t, trie.Get("800"))
	assert.Nil(t, trie.Get("810"))

	assert.Nil(t, trie.BestMatch("1"))
	assert.Nil(t, trie.BestMatch("12"))
	assert.Nil(t, trie.BestMatch("1700"))
	assert.Nil(t, trie.BestMatch("18"))
	assert.Nil(t, trie.BestMatch("180"))
	assert.Nil(t, trie.BestMatch("1802"))
	assert.Nil(t, trie.BestMatch("18001"))
	assert.Nil(t, trie.BestMatch("1812"))

	assert.Nil(t, trie.BestMatch("2"))
	assert.Nil(t, trie.BestMatch("28"))
	assert.Nil(t, trie.BestMatch("8"))
	assert.Nil(t, trie.BestMatch("80"))
	assert.Nil(t, trie.BestMatch("81"))

	trie.Put("1", "a")
	trie.Put("18", "b")
	trie.Put("1800", "c")
	trie.Put("181", "d")
	trie.Put("800", "e")
	trie.Put("810", "f")

	assert.Equal(t, tricks.ValPtr("a"), trie.Get("1"))
	assert.Equal(t, tricks.ValPtr("b"), trie.Get("18"))
	assert.Equal(t, tricks.ValPtr("c"), trie.Get("1800"))

	assert.Nil(t, trie.Get("12"))
	assert.Nil(t, trie.Get("18001"))
	assert.Nil(t, trie.Get("801"))
	assert.Nil(t, trie.Get("700"))

	assert.Equalf(t, tricks.ValPtr("a"), trie.BestMatch("1"), "got: %+v", tricks.PtrVal(trie.BestMatch("1")))
	assert.Equalf(t, tricks.ValPtr("a"), trie.BestMatch("12"), "got: %+v", tricks.PtrVal(trie.BestMatch("12")))
	assert.Equalf(t, tricks.ValPtr("a"), trie.BestMatch("1700"), "got: %+v", tricks.PtrVal(trie.BestMatch("1700")))
	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("18"), "got: %+v", tricks.PtrVal(trie.BestMatch("18")))
	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("180"), "got: %+v", tricks.PtrVal(trie.BestMatch("180")))
	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("1802"), "got: %+v", tricks.PtrVal(trie.BestMatch("1802")))
	assert.Equalf(t, tricks.ValPtr("c"), trie.BestMatch("18001"), "got: %+v", tricks.PtrVal(trie.BestMatch("18001")))
	assert.Equalf(t, tricks.ValPtr("d"), trie.BestMatch("1812"), "got: %+v", tricks.PtrVal(trie.BestMatch("1812")))
	assert.Equalf(t, tricks.ValPtr("e"), trie.BestMatch("8001"), "got: %+v", tricks.PtrVal(trie.BestMatch("8001")))
	assert.Equalf(t, tricks.ValPtr("e"), trie.BestMatch("800100"), "got: %+v", tricks.PtrVal(trie.BestMatch("800100")))
	assert.Equalf(t, tricks.ValPtr("f"), trie.BestMatch("8101"), "got: %+v", tricks.PtrVal(trie.BestMatch("8101")))
	assert.Equalf(t, tricks.ValPtr("f"), trie.BestMatch("810100"), "got: %+v", tricks.PtrVal(trie.BestMatch("810100")))

	assert.Nilf(t, trie.BestMatch("2"), "got: %+v", tricks.PtrVal(trie.BestMatch("2")))
	assert.Nilf(t, trie.BestMatch("28"), "got: %+v", tricks.PtrVal(trie.BestMatch("28")))
	assert.Nilf(t, trie.BestMatch("8"), "got: %+v", tricks.PtrVal(trie.BestMatch("8")))
	assert.Nilf(t, trie.BestMatch("80"), "got: %+v", tricks.PtrVal(trie.BestMatch("80")))
	assert.Nilf(t, trie.BestMatch("81"), "got: %+v", tricks.PtrVal(trie.BestMatch("81")))

	trie.Delete("1800")

	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("1800"), "got: %+v", tricks.PtrVal(trie.BestMatch("1800")))
	assert.Equalf(t, tricks.ValPtr("d"), trie.BestMatch("181"), "got: %+v", tricks.PtrVal(trie.BestMatch("181")))
	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("18"), "got: %+v", tricks.PtrVal(trie.BestMatch("18")))
	assert.Equalf(t, tricks.ValPtr("a"), trie.BestMatch("1"), "got: %+v", tricks.PtrVal(trie.BestMatch("1")))

	trie.Put("1800", "c2")

	assert.Equalf(t, tricks.ValPtr("c2"), trie.BestMatch("1800"), "got: %+v", tricks.PtrVal(trie.BestMatch("1800")))
	assert.Equalf(t, tricks.ValPtr("d"), trie.BestMatch("181"), "got: %+v", tricks.PtrVal(trie.BestMatch("181")))
	assert.Equalf(t, tricks.ValPtr("b"), trie.BestMatch("18"), "got: %+v", tricks.PtrVal(trie.BestMatch("18")))
	assert.Equalf(t, tricks.ValPtr("a"), trie.BestMatch("1"), "got: %+v", tricks.PtrVal(trie.BestMatch("1")))
}
