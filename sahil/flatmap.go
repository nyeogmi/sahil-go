package sahil

// FlatMap takes an existing Paginated and a function from elements to Paginated,
// then combines the two to generate a new Paginated.
//
// FlatMap(p, fn) is equivalent to Flatten(p, fn), but may be implemented more
// efficiently in practice. (Currently, it is not.)
func FlatMap[A any, B any](
	p Paginated[A],
	fn func(A) (Paginated[B], error),
) Paginated[B] {
	return Flatten(Map(p, fn))
}
