# sahil

`sahil` is a helper library for pipelines in Go.

It exists to avoid the N+1 queries problem in common database situations, without resorting to a single massive query.

Specifically, it targets the use case where:

- operating in batches is efficient
- it takes an unpredictable number of input items to generate your desired output
- operating on all your input items at once would be just too much

It gets you predictably-sized chunks of output data when you're writing code that takes a giant blob of input, then uses blocklists and searches to filter or inflate it unpredictably.

This abstraction is probably older than my library, but I independently came up with it as a part of my job and it helped us improve performance across a bunch of APIs. I thought it might be useful to other people.

## Usage

`sahil` defines one public type, called `Paginated`. `Paginated` is an iterator type that operates in batches.

There are several functions to construct a `Paginated`. If you're getting data from an external source, your options are:

- `Empty()`: constructs a `Paginated` with no elements
- `Channel(chan)`: constructs a `Paginated` which will deplete `chan` message-by-message
- `Func(f)`: constructs a `Paginated` which calls `f` every time it needs an element
- `Slice([]int {1, 2, 3})`: constructs a `Paginated` whose elements are 1, 2, 3
- `SliceFunc(f)`: constructs a `Paginated` which calls `f` to get a slice of elements every time it needs an element

Each of these comes with caveats that are explained inside the documentation.

From there, you can build new `Paginated` instances with a variety of helper functions. `sahil` provides the standard functional programming primitives:

- `Filter(pg, fn)`: takes a `Paginated` and drops all elements that fail to satisfy a condition
- `Flatten(pg)`: takes a `Paginated` of `Paginated` and strings together the results
- `FlatMap(pg, fn)`: takes the elements of a `Paginated` and calls a function on each to get a new `Paginated`, then strings them together
- `Map(pg, fn)`: takes the elements of a `Paginated` and calls a function on each

(You're encouraged not to use these more than needed, since functional code can be hard to debug.)

It provides one additional function that is unusual:

- `WindowedMap(pg, fn)`: operates on a `Paginated` in small batches -- estimating the size needed based on the ratio of input elements to output elements in previous batches

Once you have a `Paginated`, it has two methods:

- `pg.Fetch(n)`: produces between n and n * 2 elements of output, unless the data source is depleted
- `pg.FetchMany(n, m)`: produces between n and m elements of output, unless the data source is depleted

If the data source is depleted, Fetch and FetchMany will produce whatever is left, then `nil` on any future call.

## Safety warnings

It's recommended that you `Fetch()` if possible, because a `Paginated` can always elect to produce fewer than `m` elements of output, as an implementation detail. `FetchMany()` is only useful if it is completely unacceptable to receive more than a certain number of elements.

If some part of a `Paginated` produces an error, `Fetch` and `FetchMany` will produce 0 elements of output, then reproduce that error on all future calls.

## A grudging note on style

Sahil's API is written in a functional style. It does not use channels or goroutines internally, except in the `Channel` constructor. This is bad style in Go.

Unfortunately, there's not really a way to provide the API I wanted without a little FP. Because `sahil` manually estimates the size of your code's needed input and re-chunks your output into acceptably large slices, your code pretty much has to run inside a bubble where it doesn't know what's calling it or what it's calling into. The glue code it's replacing is in an awkward place where you probably want visibility into your stack but can't easily get it.

I generally don't recommend writing functional programming-flavored code in any language where other paradigms are available: it typically leads to complicated stack traces and will probably break your debugger. 

To avoid the problem of having no clue what went wrong and no ready way to debug it, I recommend treating `sahil` the way you would treat IO inside your program -- unit test the parts of your program below it if possible, and try to expose APis that don't depend on it.

## Future work

It might be good to:

- be able to interleave two Paginators
- be able to shuffle a rotating buffer of ~50 elements via Paginators
- use `sahil` in pre-generics versions of Go
- use goroutines to expose a paginator that concurrently stays a few elements ahead of its consumer
- provide more explicit support for fanning out to multiple consumers
- provide more support for very long-lived Paginators -- this currently is best done with channels
- write more unit tests

## Licensing

`sahil` is dual-licensed under the Apache-2.0 and MIT licenses. If you need something else, let me know and we can work something out.