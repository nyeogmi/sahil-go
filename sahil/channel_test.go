package sahil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	ch := make(chan string)
	go func() {
		ch <- "Hello,"
		ch <- "it's"
		ch <- "me!"
		ch <- "A cool bat."
	}()

	src := Channel(ch)

	results, err := src.Fetch(2)
	assert.Equal(t, []string{"Hello,", "it's"}, results)
	assert.Nil(t, err)

	results, err = src.Fetch(2)
	assert.Equal(t, []string{"me!", "A cool bat."}, results)
	assert.Nil(t, err)

	finished := false
	done := make(chan bool)
	defer close(done)
	go func() {
		results, err = src.Fetch(2)
		finished = true
		assert.Equal(t, 0, len(results))
		assert.Nil(t, err)
		done <- true
	}()

	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, finished, false) // not until channel is closed
	close(ch)

	select {
	case <-done:
		assert.Equal(t, finished, true)
	case <-time.After(250 * time.Millisecond):
		t.Error("should have resumed and sent 'done'")
	}
}
