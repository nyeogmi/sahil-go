package sahil

// Channel wraps a Go channel such that its messages can be fetched via the Paginated
// interface.
//
// If the channel is empty, then Fetch and FetchMany will hang until the channel
// is closed.
//
// Note that clients of the Channel can attempt to Fetch more elements than expected.
// For instance, MapWindowed may attempt to fetch 3 elements if asked for 2.
//
// Asking the Channel for more elements than there are messages will cause Fetch
// to block until the channel is closed or the messages are sent. To avoid this,
// it's recommended that you close the channel as early as possible and that you
// avoid writing your code around the assumption that an exact number of messages
// will be available.
//
// (Future versions may avoid this concurrency-related issue.)
func Channel[A any](channel chan A) Paginated[A] {
	return Func(func() (A, error) {
		a, ok := <-channel
		if !ok {
			return a, EOF
		}
		return a, nil
	})
}
