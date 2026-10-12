# Tricks [![Coverage Badge][bdg-cov-tricks]][action-tests]
Frequent magics every developer may need. This package contains general-purpose technical utilities.
Full documentation and examples are available in [pkg.go.dev][go-pkg].

## Functional options

`Option`, `MutableOption`, and `ImmutableOption` update a value one field at a time.
`ApplyOptions` applies them in order. See `ExampleApplyOptions_basic`.

## Matchers

`Matcher` answers whether a value meets a criterion. `MatchAnd`, `MatchOr`, `MatchXor`,
`MatchNot`, `MatchEqual`, `MatchGreaterThan`, and `MatchLesserThan` combine those answers.
Sequence helpers such as `SliceFilter` take a matcher.

## Sequences

`SliceMap`, `SliceFilter`, `SliceReduce`, `SliceFlat`, and `SliceDistinct` are the slice forms.
The same names without the `Slice` prefix walk an `iter.Seq`. `OrganicSliceMap` and
`OrganicSliceFilter` do the same work without allocating an iterator.

## Middleware

`MiddlewareStack.Push` chains handlers. The first middleware pushed runs first.

## mathx

`github.com/janstoon/toolbox/tricks/mathx` adds integer gcd, lcm, primes, Euler's totient,
and primitive roots.

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-tricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_tricks_refs_heads_master.json
[go-pkg]: https://pkg.go.dev/github.com/janstoon/toolbox/tricks
