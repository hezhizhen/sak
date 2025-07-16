// Package utils provides utility functions for error handling.
package utils

// CheckError checks if the provided error is not nil and panics if it is.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
