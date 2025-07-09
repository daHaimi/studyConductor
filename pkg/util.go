package pkg

func BoolP(b bool) *bool {
	return &b
}

func TypedSlice[T any](collection []any) []T {
	var typedSlice []T
	for _, concrete := range collection {
		typedSlice = append(typedSlice, concrete.(T))
	}
	return typedSlice
}
