package main

import (
	"../govec"
	"fmt"
	"net"
	"os"
)

var Logger *govec.GoLog

func main() {
	Logger = govec.Initialize("Client", "testclient.log")
	// check if the number of arguments passing to this program is correct
	if len(os.Args) != 1 {
		fmt.Println("program is supposed to have 0 arguments !\n")
	}

	rAddr, errR := net.ResolveUDPAddr("udp4", ":8080")
	printErr(errR)
	lAddr, errL := net.ResolveUDPAddr("udp4", ":18585")
	printErr(errL)

	conn, errDial := net.DialUDP("udp", lAddr, rAddr)
	printErr(errDial)

	// sending UDP packet to specified address and port
	msg := "get me the message !"
	_, errWrite := conn.Write(Logger.PrepareSend("Asking time", []byte(msg)))
	printErr(errWrite)

	// Reading the response message
	var buf [1024]byte
	n, errRead := conn.Read(buf[0:])
	printErr(errRead)
	//@dump
	fmt.Println(string(Logger.UnpackReceive("Received", buf[:n])))

	os.Exit(0)
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
