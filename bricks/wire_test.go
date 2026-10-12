package bricks_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
)

func TestCodedError(t *testing.T) {
	err := bricks.ErrorWithCode(bricks.ErrCodeNotFound, errors.New("gone"))
	var coded bricks.Coded
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, bricks.ErrCodeNotFound, coded.Code())
	assert.EqualError(t, errors.Unwrap(err), "gone")
}

func TestParseError(t *testing.T) {
	assert.NoError(t, bricks.ParseError(nil, bricks.ErrUnknown))

	err := bricks.ParseError(context.Canceled, bricks.ErrUnknown)
	require.ErrorIs(t, err, bricks.ErrCanceled)
	require.ErrorIs(t, err, context.Canceled)

	err = bricks.ParseError(context.DeadlineExceeded, bricks.ErrUnknown)
	require.ErrorIs(t, err, bricks.ErrDeadlineExceeded)

	err = bricks.ParseError(os.ErrDeadlineExceeded, bricks.ErrUnknown)
	require.ErrorIs(t, err, bricks.ErrDeadlineExceeded)

	err = bricks.ParseError(&net.DNSError{IsTimeout: true}, bricks.ErrUnknown)
	require.ErrorIs(t, err, bricks.ErrDeadlineExceeded)

	err = bricks.ParseError(errors.New("other"), bricks.ErrUnknown)
	require.ErrorIs(t, err, bricks.ErrUnknown)
	assert.EqualError(t, err, "unknown error\nother")
}

func ExampleParseError() {
	err := bricks.ParseError(context.Canceled, bricks.ErrUnknown)
	fmt.Println(errors.Is(err, bricks.ErrCanceled))
	fmt.Println(bricks.ParseError(nil, bricks.ErrUnknown))

	// Output:
	// true
	// <nil>
}

func ExampleQueue() {
	var q bricks.Queue[string]
	q.Enqueue("a", "b")
	q.Enqueue("c")

	for v := range q.All() {
		fmt.Println(v)
	}

	// Output:
	// a
	// b
	// c
}

func ExampleRegistry() {
	reg := bricks.NewRegistry[int]()
	_ = reg.Register("one", 1)
	v, err := reg.Get("one")
	fmt.Println(v, err == nil)
	_, err = reg.Get("missing")
	fmt.Println(errors.Is(err, bricks.ErrNotFound))

	// Output:
	// 1 true
	// true
}

func ExampleCompensatorAsError() {
	err := bricks.CompensatorAsError(bricks.CompensatorFunc(func(_ context.Context, err error) error {
		fmt.Println("compensate", err)

		return nil
	}), errors.New("step failed"))

	var c bricks.Compensator
	if errors.As(err, &c) {
		_ = c.Compensate(context.Background(), err)
	}

	// Output:
	// compensate step failed
}

func ExampleTrie() {
	trie := bricks.Trie[string, rune, string](func(s string) []rune { return []rune(s) })
	trie.Put("18", "area")
	trie.Put("1800", "toll-free")

	fmt.Println(*trie.Get("1800"))
	fmt.Println(*trie.BestMatch("1800555"))
	fmt.Println(trie.Get("missing") == nil)

	// Output:
	// toll-free
	// toll-free
	// true
}
