package sahil

import (
	"math"
)

type mapWindowed[A any, B any] struct {
	underlying Paginated[A]
	fn         func([]A) ([]B, error)
	nIn, nOut  int
}

// MapWindowed applies a function to a series of implementation-chosen windows of
// the Paginated. It concatenates the slices of values produced by that function
// into a new Paginated.
//
// The function provided takes an arbitrary window of unspecified size and may
// produce any number of output elements. This means it can simulate Filter and
// FlatMap.
//
// The implementation of MapWindowed assumes that operating on more elements than
// strictly necessary is slightly expensive, and operating in more _batches_ than
// strictly necessary is highly expensive.
//
// It estimates the number of input elements it will need to produce its output
// based on the number of output elements it has produced for its input so far:
// then it tries to generate its output using a minimal number of input element
// batches.
//
// The current implementation uses a bunch of heuristics that made practical
// sense at my job, but those heuristics aren't set in stone.
//
func MapWindowed[A any, B any](
	p Paginated[A],
	fn func([]A) ([]B, error),
) Paginated[B] {
	return wrap[B](&mapWindowed[A, B]{
		underlying: p,
		fn:         fn,
		nIn:        0,
		nOut:       0,
	})
}

// used when calling down to our source
const pessimismFactorSmall = 1.2 // 20% more than we think we need
const pessimismFactorBig = 1.5   // 50% more than we think we need

func (m *mapWindowed[A, B]) Fetch(atLeast int) ([]B, error) {
	var results []B

	// take at most 10 batches to get everything
	// (to avoid the results just trickling in towards the end of input)
	var smallestFetchAllowed = int(atLeast / 10)
	if smallestFetchAllowed < 1 {
		smallestFetchAllowed = 1
	}

	for {
		proportion := float64(m.nOut+1.0) / float64(m.nIn+1.0)

		optimisticInput := float64(atLeast) / proportion
		atLeastInput := pessimismFactorSmall * optimisticInput
		atMostInput := pessimismFactorBig * atLeastInput

		if atLeastInput < float64(smallestFetchAllowed) {
			atLeastInput = float64(smallestFetchAllowed)
		}

		input, err := m.underlying.FetchRange(int(atLeastInput), int(math.Ceil(atMostInput)))
		if err != nil {
			return nil, err
		}

		output, err := m.fn(input)
		if err != nil {
			return nil, err
		}

		m.nIn += len(input)
		m.nOut += len(output)
		if results == nil {
			results = output
		} else {
			results = append(results, output...)
		}

		if len(results) >= atLeast || *m.underlying.isExhausted {
			return results, nil
		}
	}
}
