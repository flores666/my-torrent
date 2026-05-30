package slicesextensions

import (
	"maps"
	"math/rand"
	"slices"
)

// PickN picks n random values from slice p
func PickN[T any](p []T, n int) []T {
	length := len(p)

	if n >= length {
		return p
	}

	picked := make(map[int]T)

	for range n {
		r := rand.Intn(length)
		if _, added := picked[r]; !added {
			picked[r] = p[r]
		}

		continue
	}

	return slices.Collect(maps.Values(picked))
}
