package sahil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	src := Empty[int]()
	results, err := src.Fetch(0)
	assert.Equal(t, len(results), 0)
	assert.Nil(t, err)
}
