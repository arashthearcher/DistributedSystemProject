// This is the package comment.
package main

import (
	"fmt"
)

const hello = "Hello, World!"

var foo = hello // line comment 2

func main() {
	//@dump
	fmt.Println(hello) // line comment 3
	fmt.Println("a")
	hello = b
	c = hello
	fmt.Println(foo) //@dump
	//@dump
}
