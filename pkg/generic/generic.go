package generic

func ToMap[K comparable, V any](list []V, f func(V) K) map[K]V {
	res := make(map[K]V, len(list))
	for _, item := range list {
		res[f(item)] = item
	}

	return res
}
