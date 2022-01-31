package sahil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	src := Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{1, 2}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(9)
	assert.EqualValues(t, []int{3, 4, 5, 6, 7, 8, 9, 10}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(9)
	assert.Equal(t, 0, len(results))
	assert.Nil(t, err)
}
