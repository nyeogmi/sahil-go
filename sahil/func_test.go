package sahil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	strings := []string{"Ik", "ben", "de", "mol"}
	i := -1

	src := Func(func() (string, error) {
		i += 1
		if i >= len(strings) {
			return "", EOF
		}
		return strings[i], nil
	})

	results, err := src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"Ik", "ben"}, results)

	results, err = src.Fetch(3)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"de", "mol"}, results)
}

func TestSliceFunc(t *testing.T) {
	strings := [][]string{{"Ik", "ben"}, {"de", "mol"}}
	i := -1

	src := SliceFunc(func() ([]string, error) {
		i += 1
		if i >= len(strings) {
			return nil, EOF
		}
		return strings[i], nil
	})

	results, err := src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"Ik", "ben"}, results)

	results, err = src.Fetch(3)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"de", "mol"}, results)
}

func TestFuncErr(t *testing.T) {
	strings := []string{"Ik", "ben", "de", "mol"}
	i := -1

	src := Func(func() (string, error) {
		i += 1
		if i >= len(strings) {
			return "", EOF
		}
		if i == 3 {
			return "", errors.New("mole error 3")
		}
		return strings[i], nil
	})

	results, err := src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"Ik", "ben"}, results)

	results, err = src.Fetch(3)
	assert.EqualError(t, err, "mole error 3")
	assert.Equal(t, 0, len(results))
}
