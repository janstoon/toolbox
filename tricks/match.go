package tricks

func MatchEqual[T comparable](expected T) func(actual T) bool {
	return func(actual T) bool {
		return expected == actual
	}
}
