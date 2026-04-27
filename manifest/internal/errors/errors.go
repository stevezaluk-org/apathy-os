package errors

import "fmt"

// Wrap Wraps error `w` with error `e`
func Wrap(w, e error) error {
	return fmt.Errorf("%w (%v)", w, e)
}
