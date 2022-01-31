package sahil

// Empty is a Paginated with no elements.
//
// Fetch() operations will immediately return a length-0 slice.
func Empty[A any]() Paginated[A] {
	return signal[A](nil)
}
