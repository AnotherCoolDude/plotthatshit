package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	defer fmt.Println("\nleaving main")
	readCSV("/Users/christianhovenbitzer/go/src/github.com/AnotherCoolDude/plotthatshit/heartbeat.csv")

	/* glot, err := glot.NewPlot(2, false, true)
	handleErr(err)
	points := [][]float64{{10, 33, 55}, {4, 6, 8}}
	glot.AddPointGroup("Heartbeat Timeline", "points", points)
	glot.SetXLabel("X-Achsis")
	glot.SetYLabel("Y-Achsis")
	glot.SetXrange(0, len(points[0]))
	glot.SetYrange(0, len(points[1]))

	glot.SavePlot("plot.png") */
}

// structs
type userCollection struct {
	userColl *[]user
}

type user struct {
	id   string
	data *[]heartBeatData
}

type heartBeatData struct {
	// maps a value to a time, e.g. heatBeat[68] = 3.56
	heartBeat map[int]int
}

// funcs

func readCSV(filename string) {

	file, err := os.Open(filename)

	if err != nil {
		fmt.Printf("error opening file %s: %s", filename, err)
	} else {
		stat, _ := file.Stat()
		fmt.Printf("using file: %v, size: %d\n\n\n", stat.Name(), stat.Size())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'

	recs := [][]string{}
	if err == nil {
		for {
			//[value, Proband, Zeit]
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				handleErr(err)
			}
			recs = append(recs, record)
		}
	}
	user := []user{&[]heartBeatData{}}
	for i, r := range recs {
		fmt.Printf("record %d: \n %v\n", i, r)

	}
}

// helper

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func newCollection() userCollection {

	user := user{
		id:   "",
		data: &[]heartBeatData{},
	}
	return userCollection{
		userColl: &[]user,
	}
}
