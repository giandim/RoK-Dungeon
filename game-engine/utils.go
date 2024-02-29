package main

import (
	"reflect"
	"syscall/js"
)

// Converts a []string to a JavaScript array
func jsSliceOf(strSlice []string) js.Value {
	jsArray := js.Global().Get("Array").New(len(strSlice))
	for i, str := range strSlice {
		jsArray.SetIndex(i, str)
	}
	return jsArray
}

// Converts a js array to a slice
func convertArrayToSlice(jsArray js.Value, convertFunc func(js.Value) interface{}) interface{} {
	length := jsArray.Length()
	resultSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(convertFunc(js.Value{}))), length, length)

	for i := 0; i < length; i++ {
		resultSlice.Index(i).Set(reflect.ValueOf(convertFunc(jsArray.Index(i))))
	}

	return resultSlice.Interface()
}

func conditionalAttribute(condition bool, attr string) string {
	if condition {
		return attr
	}
	return ""
}
