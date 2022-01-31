package sahil

import "errors"

type fetchFunc[T any] struct {
	fn func() (T, error)
}

// EOF signals end of file for Func.
var EOF = errors.New("no more elements")

// Func wraps a function such that its return values are the elements of a
// Paginated.
//
// The function is called until it produces an error of EOF.
//
// Any error other than EOF will be propagated to the Fetch or FetchMany caller.
func Func[T any](fn func() (T, error)) Paginated[T] {
	return wrap[T](&fetchFunc[T]{fn})
}

// SliceFunc wraps a slice-returning function such that the elements of its
// returned slices are the elements of a Paginated.
//
// The function is called until it produces an error of EOF.
//
// Repeatedly producing an empty slice may cause the caller to loop
// infinitely in search of an element.
func SliceFunc[T any](fn func() ([]T, error)) Paginated[T] {
	return FlatMap(Func(fn), func(ts []T) (Paginated[T], error) {
		return Slice(ts), nil
	})
}

func (f *fetchFunc[T]) Fetch(atLeast int) ([]T, error) {
	var out []T

	for len(out) < atLeast {
		t, err := f.fn()
		if errors.Is(err, EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
