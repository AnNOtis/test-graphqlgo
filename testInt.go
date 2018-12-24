package main

import (
	"fmt"
	"reflect"
)

func main() {
	var i32 int32 = 32
	var i64 int32 = 64
	var i int = 123456

	fmt.Println(reflect.TypeOf(int64(i32)))
	fmt.Println(i64)
	fmt.Println(i)
}
