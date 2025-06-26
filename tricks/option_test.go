package tricks_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestApplyOptions(t *testing.T) {
	type bag struct {
		number int
	}

	oopAdder := tricks.ImmutableOption[bag](func(s bag) bag {
		s.number++

		return s
	})

	v := bag{3}
	mv := tricks.ApplyOptions(&v, oopAdder)
	assert.NotEqual(t, v, tricks.PtrVal(mv))
	assert.NotEqual(t, &v, mv)
	assert.Equal(t, 3, v.number)
	assert.Equal(t, 4, mv.number)

	ipAdder := tricks.MutableOption[bag](func(s *bag) {
		s.number++
	})
	v = bag{3}
	mv = tricks.ApplyOptions(&v, ipAdder)
	assert.Equal(t, v, tricks.PtrVal(mv))
	assert.Equal(t, &v, mv)
	assert.Equal(t, 4, v.number)
	assert.Equal(t, 4, mv.number)
}

func ExampleApplyOptions_basic() {
	type configuration struct {
		scope         string
		verbose       bool
		skipSignature bool
		maxRetries    int
	}

	options := make([]tricks.Option[configuration], 0)
	options = append(options, tricks.ImmutableOption[configuration](func(s configuration) configuration {
		s.scope = "example"

		return s
	}))

	options = append(options, tricks.MutableOption[configuration](func(s *configuration) {
		s.verbose = true
	}))

	options = append(options, tricks.ImmutableOption[configuration](func(s configuration) configuration {
		s.maxRetries = 5

		return s
	}))

	cfg := new(configuration)
	cfg = tricks.ApplyOptions(cfg, options...)

	fmt.Println(cfg)
	// Output: &{example true false 5}
}
