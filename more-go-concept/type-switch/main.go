package main

import (
	"fmt"
	"reflect"
)

func main() {
	var t1 string = "this is a string"
	var t2 *string = &t1
	discoverType(t1)
	discoverType(t2)
	var t3 int = 1
	discoverType(t3)
	discoverType(nil)
}

func discoverType(t any) {
	switch v := t.(type) {
	case string:
		t2 := v + " and this is another string"
		fmt.Printf("String found: %s\n", t2)
	case *string:
		fmt.Printf("Pointer to string found: %s\n", *v)
	default:
		// fmt.Printf("Unknown type: %T\n", v)
		// fmt.Printf("Unknown type: %s\n", reflect.TypeOf(v))
		myType := reflect.TypeOf(v)
		if myType == nil {
			fmt.Printf("Unknown type: nil\n")
		} else {
			fmt.Printf("Unknown type: %s\n", myType)
		}
	}
}
