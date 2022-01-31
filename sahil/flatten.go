package sahil

type flatten[T any] struct {
	source    Paginated[Paginated[T]]
	buf       []Paginated[T]
	exhausted bool
	pagErr    error
}

// Concat takes a slice of Paginated and combines them into a single Paginated.
//
// Concat(ps...) is equivalent to Flatten(Slice(ps)), but may be implemented more
// efficiently in practice. (Currently, it is not.)
func Concat[T any](ps ...Paginated[T]) Paginated[T] {
	return Flatten(Slice(ps))
}

// Flatten takes a Paginated of Paginated and combines them into a single Paginated.
//
// It does this by fetching each Paginated individually from the source, then
// yielding each of its elements in sequence.
//
// For instance, Flatten(Slice([]Paginated[int] { Slice(1, 2, 3), Slice(4, 5, 6) }))
// is equivalent to Slice(1, 2, 3, 4, 5, 6).
func Flatten[T any](p Paginated[Paginated[T]]) Paginated[T] {
	return wrap[T](&flatten[T]{
		source: p,
		buf:    nil,
	})
}

func (j *flatten[T]) Fetch(atLeast int) ([]T, error) {
	var out []T

	for {
		if len(out) >= atLeast {
			return out, nil
		}

		current, err := j.currentPaginated()
		if err != nil {
			return nil, err
		}

		if current == nil {
			return out, nil
		}

		buf, err := current.Fetch(atLeast - len(out))
		if err != nil {
			return nil, err
		}

		if out == nil {
			out = buf
		} else {
			out = append(out, buf...)
		}
	}
}

func (j *flatten[T]) currentPaginated() (*Paginated[T], error) {
	if j.pagErr != nil {
		return nil, j.pagErr
	}

	if j.exhausted {
		return nil, nil
	}

	for {
		if len(j.buf) == 0 {
			nextPaginators, err := j.source.Fetch(1)

			if err == nil {
				j.buf = nextPaginators
			} else {
				j.buf = nil
				j.pagErr = err
				return nil, j.pagErr
			}
		}

		for len(j.buf) > 0 {
			next := j.buf[0]

			if *next.isExhausted {
				j.buf = j.buf[1:]
			}

			return &next, nil
		}

		if *j.source.isExhausted {
			// Shortcut to this path so we can get here without trying and failing to find another one
			j.exhausted = true
			return nil, nil
		}
	}
}
