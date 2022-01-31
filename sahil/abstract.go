package sahil

import (
	"sync"
)

// fetch is the internal interface implemented by Paginated-compatible data
// sources. It can return as many results as it wants, but must return at least
// atLeast results so long as the  data sourec is not completed.
type fetch[T any] interface {
	Fetch(atLeast int) ([]T, error)
}

// Paginated is a struct for retrieving elements in batches from a data source.
//
// It's designed to replace channel pipelines in programs where operating element-by-
// element is inefficient, but operating on everything at once would do redundant work.
//
// If you know how many output items you want but not how many input items you will
// need to look at to get that much output, Paginated can estimate that, then attempt
// to do the work in batches.
//
// Its methods are Fetch(atLeast), which returns at least atLeast elements, and
// FetchRange(atLeast, atMost), which returns at least atLeast elements and at most
// atMost elements. (inclusive)
//
// Each method will produce less than `atLeast` elements once the data source runs out.
// Further calls will produce nil.
//
// Any error will result in the immediate end of output. (In other words, you can't
// recover any output generated before the error occurred.) Future calls to
// Fetch or FetchMany will produce that error.
//
// Paginated is thread-safe.
type Paginated[T any] struct {
	underlying   buffered[T]
	mutex        *sync.Mutex
	isExhausted  *bool // pointer so it isn't inadvertently copied
	err          *error
	atMostFactor float64
}

// buffered augments Paginated with buffering behavior -- if a fetch implementor
// returns too many results, buffered will temporarily store all elements after
// the first atMost elements.
type buffered[T any] struct {
	underlying *fetch[T]
	buffer     *[]T
}

// wrap converts a fetch (the internal interface) to a Paginated (the external one)
func wrap[T any](f fetch[T]) Paginated[T] {
	isExhausted := false
	var err error
	var buf []T
	return Paginated[T]{
		underlying: buffered[T]{
			underlying: &f,
			buffer:     &buf,
		},
		mutex:        &sync.Mutex{},
		isExhausted:  &isExhausted,
		err:          &err,
		atMostFactor: 2.0,
	}

}

// signal generates a paginated from an error code.
//
// Not providing that error code will result in an empty Paginated.
func signal[A any](err error) Paginated[A] {
	exh := true
	return Paginated[A]{
		underlying: buffered[A]{
			underlying: nil,
			buffer:     nil,
		},
		mutex:        &sync.Mutex{},
		isExhausted:  &exh,
		err:          &err,
		atMostFactor: 2.0,
	}
}

// Fetch fetches at least `atLeast` elements from the underlying `Fetch` object.
// If there are not `atLeast` elements, it will reproduce whatever was found.
//
// Further calls (or errors) result in nil. Errors additionally result in an
// error value.
//
// Equivalent to FetchRange(atLeast, atLeast * 2).
func (p Paginated[T]) Fetch(atLeast int) ([]T, error) {
	return p._fetch(atLeast, int(atLeast*int(p.atMostFactor)))
}

// FetchRange fetches at least `atLeast` elements from the underlying `Fetch`
// object. If there are not `atLeast` elements, it will reproduce whatever was
// found.
//
// Further calls (or errors) result in nil. Errors additionally result in an
// error value.
func (p Paginated[T]) FetchRange(atLeast, atMost int) ([]T, error) {
	// if it overrides the atMostFactor, respect that override
	atMost2 := int(float64(atLeast) * p.atMostFactor)
	if atMost > atMost2 {
		atMost = atMost2
	}
	if atMost < atLeast {
		atMost = atLeast
	}
	return p._fetch(atLeast, atMost)
}

func (p Paginated[T]) _fetch(atLeast, atMost int) ([]T, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if *p.isExhausted {
		return nil, *p.err
	}
	if atLeast == 0 {
		return nil, nil
	}

	result, err := p.underlying._fetch(atLeast, atMost)
	if len(result) < atLeast || err != nil {
		*p.isExhausted = true
		*p.underlying.underlying = nil // allow this stuff to be freed
		*p.underlying.buffer = nil
		*p.err = err
	}

	if err != nil {
		// don't produce results in this case
		result = nil
	}

	return result, err
}

func (b buffered[T]) _fetch(atLeast int, atMost int) ([]T, error) {
	if len(*b.buffer) > atMost {
		chunk := (*b.buffer)[:atMost]
		*b.buffer = (*b.buffer)[atMost:]
		return chunk, nil
	}

	if len(*b.buffer) >= atLeast {
		chunk := *b.buffer
		*b.buffer = nil
		return chunk, nil
	}

	nWanted := atLeast - len(*b.buffer)
	buf, err := (*b.underlying).Fetch(nWanted)
	if err != nil {
		return nil, err
	}
	*b.buffer = append(*b.buffer, buf...)

	// now repeat the original logic for dealing with at most
	if len(*b.buffer) > atMost {
		chunk := (*b.buffer)[:atMost]
		*b.buffer = (*b.buffer)[atMost:]
		return chunk, nil
	}

	chunk := *b.buffer
	*b.buffer = nil
	return chunk, nil
}
