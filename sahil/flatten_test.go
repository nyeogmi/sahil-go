package sahil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	src := Flatten(Slice([]Paginated[int]{
		Slice([]int{1, 2, 3}),
		Slice([]int{4, 5}),
		Empty[int](),
		Slice([]int{6, 7, 8, 9}),
	}))

	results, err := src.Fetch(3)
	assert.EqualValues(t, []int{1, 2, 3}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(4)
	assert.EqualValues(t, []int{4, 5, 6, 7}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(2)
	assert.EqualValues(t, []int{8, 9}, results)
	assert.Nil(t, err)
}

func TestConcat(t *testing.T) {
	src := Concat(
		Slice([]string{
			"Desmodus",
			"Diaemus",
		}),
		Slice([]string{
			"Diphylla",
		}),
	)

	results, err := src.Fetch(4)
	assert.EqualValues(t, []string{"Desmodus", "Diaemus", "Diphylla"}, results)
	assert.Nil(t, err)
}

func TestConcatErr(t *testing.T) {
	src := Concat(
		Slice([]string{
			"Desmodus",
			"Diaemus",
		}),
		Slice([]string{
			"Diphylla",
		}),
		Func(func() (string, error) {
			return "", errors.New("BAT ERROR")
		}),
	)

	results, err := src.Fetch(4)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "BAT ERROR")
}
