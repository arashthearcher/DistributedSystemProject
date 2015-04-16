package main

import ()

func main() {
	n := 5
	for i := 0; i < n; i++ {
		send("Q", i)
	}
	commit := decide()
	for i := 0; i < n; i++ {
		recv(buf, i)
	}
	if buf == "A" {
		commit := "A"
	}
	for i := 0; i < n; i++ {
		send(commit, i)
	}
}

func send(q string, id int) {

}

func decide() {
	result = ""
	return result
}
