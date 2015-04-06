// This is the package comment.
package main

import (
	"fmt"
)

// This comment is associated with the hello constant.
const hello = "Hello, World!" // line comment 1

// This comment is associated with the foo variable.
var foo = hello // line comment 2

//@dump This comment is associated with the main function.
func main() {
	//@dump
	fmt.Println(hello) // line comment 3
	fmt.Println("a")
	hello = b
	c = hello
	fmt.Println(foo) //@dump
	//@dump
}
