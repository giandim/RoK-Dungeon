package main

import "syscall/js"

// Converts a []string to a JavaScript array
func jsSliceOf(strSlice []string) js.Value {
	jsArray := js.Global().Get("Array").New(len(strSlice))
	for i, str := range strSlice {
		jsArray.SetIndex(i, str)
	}
	return jsArray
}
