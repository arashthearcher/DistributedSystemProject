package main

import "fmt"

var gx = "Goodbye"

func main() {
	x := "Hello"
	fmt.Println(x)
	fmt.Println(gx)
	for {
		var amt = 1
		fmt.Println("inside incr")
		return x + amt
	}
	fmt.Println(inc(2))
	f()

}

func f() {
	y := "Rocky"
	fmt.Println(y)
	fmt.Println(gx)
}
