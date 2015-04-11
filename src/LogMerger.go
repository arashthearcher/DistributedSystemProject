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

	mergedPoints := make([]Point, 0)
	for i := 0; i < len(log1); i++ {

		matchedPoints := findMatch(log1[i], log2)
		for j := 0; j < len(matchedPoints); j++ {
			mergedPoints = append(mergedPoints, mergePoints(matchedPoints[j], log1[i]))
		}
	}

	return mergedPoints

}

func mergePoints(p1, p2 Point) Point {
	var mergedPoint Point
	mergedPoint.Dump = append(p1.Dump, p2.Dump...)
	mergedPoint.LineNumber = p1.LineNumber + "-" + p2.LineNumber
	pVClock1, err := vclock.FromBytes(p1.vectorClock)
	printErr(err)
	pVClock2, err2 := vclock.FromBytes(p2.vectorClock)
	printErr(err2)
	temp := pVClock1.Copy()
	temp.Merge(pVClock2)
	mergedPoint.vectorClock = temp.Bytes()

	return mergedPoint

}

func findMatch(point Point, log []Point) []Point {
	matched := make([]Point, 0)
	pVClock, err := vclock.FromBytes(point.vectorClock)
	printErr(err)
	for i := 0; i < len(log); i++ {

		otherVClock, err2 := vclock.FromBytes(log[i].vectorClock)
		printErr(err2)
		if pVClock.Matches(otherVClock) {
			matched = append(matched, log[i])
		}
	}

	return matched
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
		fmt.Println(err)
	}
}

type Point struct {
	Dump        []NameValuePair
	LineNumber  string
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
