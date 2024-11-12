package generic

import "reflect"

func ToMap[K comparable, V any](list []V, f func(V) K) map[K]V {
	res := make(map[K]V, len(list))
	for _, item := range list {
		res[f(item)] = item
	}

	return res
}

func Map[S any, T any](list []S, f func(S) T) []T {
	res := make([]T, 0, len(list))
	for _, item := range list {
		res = append(res, f(item))
	}

	return res
}

func structToArray(s interface{}) []interface{} {
	v := reflect.ValueOf(s)
	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	return values
}
