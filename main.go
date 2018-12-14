package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/Arafatk/glot"
)

func main() {
	defer fmt.Println("\nleaving main")
	col := readCSV("/Users/christianhovenbitzer/go/src/github.com/AnotherCoolDude/plotthatshit/heartbeat.csv")

	xValues := []float64{}
	yValues := []int{}
	for _, u := range *col.getID("P01").data {
		xValues = append(xValues, u.time)
		yValues = append(yValues, u.value)
	}

	glot, err := glot.NewPlot(2, false, true)
	handleErr(err)
	points := [][]float64{{10, 33, 55}, {4, 6, 8}}
	glot.AddPointGroup("Heartbeat Timeline", "linepoints", points)
	glot.SetXLabel("X-Achsis")
	glot.SetYLabel("Y-Achsis")
	glot.SetXrange(0, int(math.Round(maxFloatSlice(xValues))))
	glot.SetYrange(0, maxIntSlice(yValues))

	glot.SavePlot("plot.png")
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
	value int
	time  float64
}

// funcs

func readCSV(filename string) userCollection {

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

	col := userCollection{
		userColl: &[]user{},
	}

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
	for _, r := range recs {
		if !col.contains(r[1]) {
			newU := user{
				id:   r[1],
				data: &[]heartBeatData{},
			}
			*col.userColl = append(*col.userColl, newU)
		}
		v, _ := strconv.Atoi(r[0])
		t, _ := strconv.ParseFloat(r[2], 64)
		newB := heartBeatData{
			value: v,
			time:  t,
		}
		*col.getID(r[1]).data = append(*col.getID(r[1]).data, newB)

		//fmt.Printf("record %d: \n %v\n", i, r)

	}
	for _, u := range *col.userColl {
		fmt.Printf("userid: %s, dataamount: %d\n", u.id, len(*u.data))
	}

	return col

}

// helper

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (col *userCollection) contains(id string) bool {
	for _, u := range *col.userColl {
		if u.id == id {
			return true
		}
	}
	return false
}

func (col *userCollection) getID(id string) *user {
	for _, u := range *col.userColl {
		if u.id == id {
			return &u
		}
	}
	return nil
}

func maxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}

func maxFloatSlice(v []float64) float64 {
	sort.Float64s(v)
	return v[len(v)-1]
}
