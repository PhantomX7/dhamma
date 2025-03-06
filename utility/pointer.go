package utility

// PointOf returns a pointer to the given value.
//
// Example:
//
//	foo := "bar"
//	fooPtr := PointOf(foo)
//	fmt.Println(*fooPtr) // Output: bar
func PointOf[T any](value T) *T {
	return &value
}
