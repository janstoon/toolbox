package bricks_test

import (
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/bricks"
)

func TestTreeRetrieval(t *testing.T) {
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
}
