package sahil

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatMap(t *testing.T) {
	src := FlatMap(
		Slice(
			[]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			},
		),
		func(s string) (Paginated[string], error) {
			words := strings.Split(s, " ")
			return Slice(words), nil
		},
	)

	results, err := src.Fetch(3)
	assert.EqualValues(t, []string{"Desmodus", "rotundus", "Diaemus"}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(3)
	assert.EqualValues(t, []string{"youngi", "Diphylla", "ecaudata"}, results)
	assert.Nil(t, err)
}

func TestFlatMapErr1(t *testing.T) {
	src := FlatMap(
		Slice(
			[]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			},
		),
		func(s string) (Paginated[string], error) {
			if strings.Contains(s, "ecaudata") {
				return Empty[string](), errors.New("BAT ERROR")
			}
			words := strings.Split(s, " ")
			return Slice(words), nil
		},
	)

	results, err := src.Fetch(1)
	assert.EqualValues(t, []string{"Desmodus"}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(4)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "BAT ERROR")
}

func TestFlatMapErr2(t *testing.T) {
	src := FlatMap(
		Concat(
			Slice([]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			}),
			Func(func() (string, error) {
				return "", errors.New("BAT ERROR")
			}),
		),
		func(s string) (Paginated[string], error) {
			words := strings.Split(s, " ")
			return Slice(words), nil
		},
	)

	results, err := src.Fetch(1)
	assert.EqualValues(t, []string{"Desmodus"}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(7)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "BAT ERROR")
}

func TestFlatMapErr3(t *testing.T) {
	src := FlatMap(
		Concat(
			Slice([]string{
				"Desmodus rotundus",
				"Diaemus youngi",
				"Diphylla ecaudata",
			}),
		),
		func(s string) (Paginated[string], error) {
			if strings.Contains(s, "Diaemus") {
				return Func(func() (string, error) {
					return "", errors.New("BAT ERROR")
				}), nil
			}
			words := strings.Split(s, " ")
			return Slice(words), nil
		},
	)

	results, err := src.Fetch(1)
	assert.EqualValues(t, []string{"Desmodus"}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(7)
	assert.Equal(t, 0, len(results))
	assert.EqualError(t, err, "BAT ERROR")
}
