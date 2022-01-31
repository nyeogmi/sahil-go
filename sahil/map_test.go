package sahil

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	src := Map(
		Slice(
			[]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			},
		),
		func(s string) (int, error) {
			return len(s), nil
		},
	)

	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{
		len("Desmodus rotundus"),
		len("Diaemus youngi"),
	}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(1)
	assert.EqualValues(t, []int{
		len("Diphylla ecaudata"),
	}, results)
	assert.Nil(t, err)
}

func TestMapErr(t *testing.T) {
	src := Map(
		Slice(
			[]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			},
		),
		func(s string) (int, error) {
			if strings.Contains(s, "ecaudata") {
				return 0, errors.New("BAT ERROR")
			}
			return len(s), nil
		},
	)

	results, err := src.Fetch(2)
	assert.EqualValues(t, []int{
		len("Desmodus rotundus"),
		len("Diaemus youngi"),
	}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(1)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "BAT ERROR")
}
