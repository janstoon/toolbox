package tricks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/tricks"
)

func TestApplyOptions(t *testing.T) {
	type bag struct {
		number int
	}

	oopAdder := tricks.OutOfPlaceOption[bag](func(s bag) bag {
		s.number++

		return s
	})

	v := bag{3}
	mv := tricks.ApplyOptions(&v, oopAdder)
	assert.NotEqual(t, v, tricks.PtrVal(mv))
	assert.NotEqual(t, &v, mv)
	assert.Equal(t, 3, v.number)
	assert.Equal(t, 4, mv.number)

	ipAdder := tricks.InPlaceOption[bag](func(s *bag) {
		s.number++
	})
	v = bag{3}
	mv = tricks.ApplyOptions(&v, ipAdder)
	assert.Equal(t, v, tricks.PtrVal(mv))
	assert.Equal(t, &v, mv)
	assert.Equal(t, 4, v.number)
	assert.Equal(t, 4, mv.number)
}
