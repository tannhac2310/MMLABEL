package generic

func ToMap[K comparable, V any](list []V, f func(V) K) map[K]V {
	res := make(map[K]V, len(list))
	for _, item := range list {
		res[f(item)] = item
	}

	return res
}

func Map[S any, T any](list []S, f func(S) T) []T {
	res := make([]T, 0, len(list))
	for i, item := range list {
		res[i] = f(item)
	}

	return res
}
