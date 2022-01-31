package sahil

type fetchSlice[T any] struct {
	slice []T
}

// Slice takes a slice of values and creates a Paginated whose values are the
// elements of that slice.
//
// You rarely need this by itself, but you can use it to provide input to other
// Paginated-consuming APIs.
func Slice[T any](slice []T) Paginated[T] {
	mySlice := make([]T, len(slice))
	copy(mySlice, slice)

	// use a lower atMostFactor to force the buffering implementation to give
	// exactly atLeast results
	out := wrap[T](&fetchSlice[T]{slice: mySlice})
	out.atMostFactor = 1.0
	return out
}

func (f *fetchSlice[T]) Fetch(atLeast int) ([]T, error) {
	// defer to the buffering implementation
	// and make sure future calls can only produce nil
	out := f.slice
	f.slice = nil
	return out, nil
}
