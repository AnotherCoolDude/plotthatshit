package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"gonum.org/v1/plot/plotter"

	"gonum.org/v1/plot"
)

func main() {
	defer fmt.Println("\nleaving main")
	col := readCSV("/Users/christianhovenbitzer/go/src/github.com/AnotherCoolDude/plotthatshit/heartbeat.csv")

	plot, err := plot.New()
	if err != nil {
		fmt.Println(err)
	}

	userMap := col.getUserMap()
	xys := make(plotter.XYs, len(userMap["P01"][0]))
	for user, coords := range userMap {
		if user == "" || user == "Proband" {
			continue
		}

		fmt.Println(user)
		for i := range coords[0] {
			fmt.Printf("X Value: %f02, Y Value: %f02\n", coords[0][i], coords[1][i])
			xys[i].X = coords[0][i]
			xys[i].Y = coords[1][i]
		}
		print(xys)
	}
	line, err := plotter.NewLine(xys)
	if err != nil {
		fmt.Println(err)
	}
	plot.Add(line)

	savePlot(plot)
	//line, err := plotter.NewLine()

	/* glot, err := glot.NewPlot(2, false, true)
	handleErr(err)

	glot.SetTitle("Auswertung")
	glot.SetXLabel("Zeit in Sekunden [s]")
	glot.SetYLabel("Herzschlag in Schlägen pro Minute [bpm]")
	glot.SetXrange(0, 250)
	glot.SetYrange(40, 100)

	for _, p := range *col.userColl {
		if p.id == "" || p.id == "Proband" {
			continue
		}
		xValues := []float64{}
		yValues := []float64{}
		for _, u := range *col.getID(p.id).data {
			xValues = append(xValues, u.time)
			yValues = append(yValues, float64(u.value))
		}
		data := [][]float64{
			xValues,
			yValues,
		}
		glot.AddPointGroup(p.id, "lines", data)
	}

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

func savePlot(plot *plot.Plot) {
	wt, err := plot.WriterTo(512, 512, "png")
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = wt.WriteTo(f)
	if err != nil {
		fmt.Println(err)
	}
}

func (col *userCollection) getUserMap() map[string][][]float64 {
	uMap := map[string][][]float64{}

	for _, p := range *col.userColl {
		if p.id == "" || p.id == "Proband" {
			continue
		}
		xValues := []float64{}
		yValues := []float64{}
		for _, u := range *col.getID(p.id).data {
			xValues = append(xValues, u.time)
			yValues = append(yValues, float64(u.value))
		}
		data := [][]float64{
			xValues,
			yValues,
		}
		uMap[p.id] = data
	}
	return uMap
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
