// LogMerger
package main

import (
	"./vclock"
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	log1 := readLog("log1.txt")
	log2 := readLog("log2.txt")
	mergedLog := mergeLogs(log1, log2)

}

func mergeLogs(log1, log2 []Point) []Point {

}

func readLog(filePath string) []Point {
	fileR, err := os.Open(filePath)
	printErr(err)

	fmt.Println("decoding " + filePath)
	decoder := gob.NewDecoder(fileR)

	pointArray := make([]Point, 0)

	var e error = nil
	for e == nil {
		var decodedPoint Point
		e = decoder.Decode(&decodedPoint)
		if e == nil {
			pointArray = append(pointArray, decodedPoint)
		}
	}

	return pointArray
}

func printErr(err error) {
	if err != nil {
		log.Println(err)
	}
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
