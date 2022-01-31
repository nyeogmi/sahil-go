package sahil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	src := Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	src = Filter(src, func(x int) (bool, error) { return x%2 == 0, nil })
	results, err := src.Fetch(2)

	// the exact number of results
	// depends on how MapWindowed estimates how many it needs
	// which is implementation-defined
	// so update this assertion based on the practical results
	assert.EqualValues(t, []int{2, 4}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(3)
	assert.EqualValues(t, []int{6, 8, 10}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(3)
	assert.Equal(t, 0, len(results))
	assert.Nil(t, err)
}

func TestFilterErr(t *testing.T) {
	src := Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	src = Filter(src, func(x int) (bool, error) {
		if x == 7 {
			return false, errors.New("LUCKY NUMBER 7")
		}
		return x%2 == 0, nil
	})

	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{2, 4}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(3)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "LUCKY NUMBER 7")
}
