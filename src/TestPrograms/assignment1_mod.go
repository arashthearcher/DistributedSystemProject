package main

import (
	"../govec"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"reflect"
)

func main() {
	InstrumenterInit()
	Logger = govec.Initialize("Client", "testclient.log")

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
	vars30 := []interface{}{conn, buf, errRead, rAddr, errDial, msg, Logger, errWrite, errR, errL, lAddr}
	varsName30 := []string{"conn", "buf", "errRead", "rAddr", "errDial", "msg", "Logger", "errWrite", "errR", "errL", "lAddr"}
	point30 := createPoint(vars30, varsName30, 30)
	encoder.Encode(point30)
	fmt.Println(string(Logger.UnpackReceive("Received", buf[:n])))

	os.Exit(0)
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var Logger *govec.GoLog

var encoder *gob.Encoder

func InstrumenterInit() {
	fileW, _ := os.Create("log_client.txt")
	encoder = gob.NewEncoder(fileW)
}

func createPoint(vars []interface{}, varNames []string, lineNumber int) Point {

	length := len(varNames)
	dumps := make([]NameValuePair, 0)
	for i := 0; i < length; i++ {

		if vars[i] != nil && ((reflect.TypeOf(vars[i]).Kind() == reflect.String) || (reflect.TypeOf(vars[i]).Kind() == reflect.Int)) {
			var dump NameValuePair
			dump.VarName = varNames[i]
			dump.Value = vars[i]
			dump.Type = reflect.TypeOf(vars[i]).String()
			dumps = append(dumps, dump)
		}
	}
	fmt.Println(Logger.GetCurrentVC())
	point := Point{dumps, string(lineNumber), Logger.GetCurrentVC()}
	fmt.Println(point.VectorClock)

	return point
}

type Point struct {
	Dump        []NameValuePair
	LineNumber  string
	VectorClock []byte
}

type NameValuePair struct {
	VarName string
	Value   interface{}
	Type    string
}

func (nvp NameValuePair) String() string {
	return fmt.Sprintf("(%s,%s,%s)", nvp.VarName, nvp.Value, nvp.Type)
}

func (p Point) String() string {
	return fmt.Sprintf("%d : %s", p.LineNumber, p.Dump)
}
