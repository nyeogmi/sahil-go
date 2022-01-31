package sahil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	if 1 == 2 {
		t.Errorf("environment is misconfigured")
	}
}

func TestNoCallsOnceExhausted(t *testing.T) {
	i := 0

	called := false
	src := Func(func() (int, error) {
		called = true
		i += 1
		if i > 5 {
			return 0, EOF
		}
		return i, nil
	})

	called = false
	results, err := src.Fetch(6)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4, 5}, results)
	assert.True(t, called)
	assert.True(t, *src.isExhausted)

	called = false
	results, err = src.Fetch(5)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(results))
	assert.False(t, called)
}

func TestFetchAtLeast(t *testing.T) {
	// do the minimum number of function calls to get the desired number of results
	i := 0
	src := Func(func() (int, error) {
		i += 1
		if i > 5 {
			return 0, EOF
		}
		return i, nil
	})

	results, err := src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2}, results)
	assert.False(t, *src.isExhausted)

	results, err = src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{3, 4}, results)
	assert.False(t, *src.isExhausted)

	results, err = src.Fetch(3)
	assert.Nil(t, err)
	assert.EqualValues(t, results, []int{5})
	assert.True(t, *src.isExhausted)
}

type bigResultsTest struct{}

func (bigResultsTest) Fetch(atLeast int) ([]int, error) {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil
}

func TestBufferBigResults(t *testing.T) {
	src := wrap[int](bigResultsTest{})

	results, err := src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4}, results)
	assert.False(t, *src.isExhausted)

	results, err = src.Fetch(2)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{5, 6, 7, 8}, results)
	assert.False(t, *src.isExhausted)

	results, err = src.Fetch(3)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{9, 10, 1, 2, 3, 4}, results)
	assert.False(t, *src.isExhausted)
}

func TestSignalEmpty(t *testing.T) {
	src := signal[int](nil)
	results, err := src.Fetch(0)
	assert.Equal(t, len(results), 0)
	assert.Nil(t, err)
}

func TestSignalErr(t *testing.T) {
	src := signal[int](errors.New("barf"))
	results, err := src.Fetch(0)
	assert.Equal(t, len(results), 0)
	assert.EqualError(t, err, "barf")
}
