package sahil

type mapFn[A any, B any] struct {
	underlying Paginated[A]
	fn         func(A) (B, error)
}

// Map applies a function to the elements in a Paginated.
//
// This results in a new Paginated with the same number of elements.
func Map[A any, B any](p Paginated[A], fn func(A) (B, error)) Paginated[B] {
	return wrap[B](&mapFn[A, B]{underlying: p, fn: fn})
}

func (m *mapFn[A, B]) Fetch(atLeast int) ([]B, error) {
	outA, err := m.underlying.Fetch(atLeast)
	if err != nil {
		return nil, err
	}

	outB := make([]B, len(outA))
	for i, a := range outA {
		b, err := m.fn(a)
		if err != nil {
			return nil, err
		}
		outB[i] = b
	}
	return outB, nil
}
