package main

import ()

func main() {
	
recv(buf)
	if buf == "Q"{
	vote:= decide()
	send(vote)
	}
if vote == "A"{
	commit = "A"
	}
else{
	recv(buf)
	commit = buf
	}
}

func send(q string, id int) {

}

func decide() {
	result = ""
	return result
}