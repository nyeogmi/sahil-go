package sahil

// Filter takes an existing paginated and applies a predicate to its elements.
//
// Elements for which the function returns true will appear in the output
// while elements that get false will not.
//
// Filter will Fetch larger batches from the underlying Paginated based on the
// proportion of items admitted early on. For instance, a Filter that rejects
// 4/5 elements early on will estimate that it needs five elements of input
// to generate each element of output.
//
// See MapWindowed for more documentation on this behavior.
func Filter[A any](
	p Paginated[A],
	fn func(A) (bool, error),
) Paginated[A] {
	return MapWindowed(p, func(as []A) ([]A, error) {
		var out []A
		for _, a := range as {
			inc, err := fn(a)
			if err != nil {
				return nil, err
			}
			if inc {
				out = append(out, a)
			}
		}
		return out, nil
	})
}
