package main

import (
	//"./govec"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
)

var encoder gob.Encoder

func InstrumenterInit() {
	fileW, _ := os.Create("log.txt")
	encoder = gob.NewEncoder(fileW)
}

var logger *govec.GoLog

func createPoint(vars []interface{}, varNames []string, lineNumber int) Point {

	length := len(varNames)
	dumps := make([]NameValuePair, length)
	for i := 0; i < length; i++ {
		dumps[i].VarName = varNames[i]
		dumps[i].Value = vars[i]
		dumps[i].Type = reflect.TypeOf(vars[i]).String()
	}
	point := Point{dumps, lineNumber, logger.currentVC}

	return point
}

type Point struct {
	Dump        []NameValuePair
	LineNumber  int
	vectorClock []byte
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

func main() {
	conn, err := net.ListenPacket("udp", ":8080")
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	printErr(err)

	for {
		if err != nil {
			printErr(err)
			continue
		}
		handleConn(conn)
		fmt.Println("some one connected!")
		vars25 := []interface{}{err, conn}
		varsName25 := []string{"err", "conn"}
		point25 := createPoint(vars25, varsName25, 25)
		encoder.Encode(point25)
	}
	conn.Close()

}

func handleConn(conn net.PacketConn) {
	var buf [512]byte

	_, addr, err := conn.ReadFrom(buf[0:])
	printErr(err)
	msg := fmt.Sprintf("Hello There! time now is %s \n", time.Now().String())
	conn.WriteTo([]byte(msg), addr)
	vars38 := []interface{}{msg, addr, conn, buf, err}
	varsName38 := []string{"msg", "addr", "conn", "buf", "err"}
	point38 := createPoint(vars38, varsName38, 38)
	encoder.Encode(point38)
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
