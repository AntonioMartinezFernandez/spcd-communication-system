package utils

import "reflect"

func Ptr[T any](v T) *T {
	return &v
}

func Val[T any](v *T) T {
	return *v
}

func InterfacePointerIsNil(val interface{}) bool {
	return reflect.ValueOf(val).IsNil()
}
