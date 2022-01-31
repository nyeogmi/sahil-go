package sahil

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapWindowed(t *testing.T) {
	src := MapWindowed(
		Slice(
			[]string{
				"A", "BB", "CCC",
				"DDDD", "EEEEE", "FFFFFF",
			},
		),
		func(strings []string) ([]int, error) {
			var results = make([]int, len(strings)*2)
			for i, s := range strings {
				results[2*i] = len(s)
				results[2*i+1] = len(s)
			}
			return results, nil
		},
	)

	// note: exact number of results depends on estimation method used by
	// MapWindowed
	// which is impl defined
	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{
		1, 1, 2, 2,
	}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(2)
	assert.EqualValues(t, []int{
		3, 3,
	}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(2)
	assert.EqualValues(t, []int{
		4, 4,
	}, results)
	assert.Nil(t, err)
}

func TestMapWindowedErr(t *testing.T) {
	src := MapWindowed(
		Slice(
			[]string{
				"A", "BB", "CCC",
				"DDDD", "EEEEE", "FFFFFF",
			},
		),
		func(ss []string) ([]int, error) {
			var results = make([]int, len(ss)*2)
			for i, s := range ss {
				results[2*i] = len(s)
				results[2*i+1] = len(s)
				if strings.Contains(s, "D") {
					return nil, errors.New("got the D")
				}
			}
			return results, nil
		},
	)

	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{
		1, 1, 2, 2,
	}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(4)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "got the D")
}
