package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gonum.org/v1/plot/vg"

	"github.com/AvraamMavridis/randomcolor"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"

	"gonum.org/v1/plot"
)

var (
	stimuli = map[int]stimulus{
		0: stimulus{
			xValue: 0.0,
			name:   "Schwarz",
		},
		1: stimulus{
			xValue: 20.0,
			name:   "Katze",
		},
		2: stimulus{
			xValue: 40.0,
			name:   "Schwarz",
		},
		3: stimulus{
			xValue: 55.0,
			name:   "Grau",
		},
		4: stimulus{
			xValue: 75.0,
			name:   "Schwarz",
		},
		5: stimulus{
			xValue: 90.0,
			name:   "RayBan",
		},
		6: stimulus{
			xValue: 110.0,
			name:   "Schwarz",
		},
		7: stimulus{
			xValue: 125.0,
			name:   "Video",
		},
		8: stimulus{
			xValue: 265.0,
			name:   "Schwarz",
		},
	}
	lPlacement = legendPlacement{
		top:     true,
		left:    true,
		xOffset: 120,
	}

	sFont, _       = vg.MakeFont("Courier", vg.Points(10))
	lFont, _       = vg.MakeFont("Courier", vg.Points(14))
	smallDrawStyle = draw.TextStyle{Font: sFont}
	largeDrawStyle = draw.TextStyle{Font: lFont}

	timeLimit   = 128.0
	userToPrint = []int{1, 2, 3, 4, 5, 7, 9, 10, 11, 12, 13}
)

const (
	imgWidth       = 2048
	imgHeight      = 1024
	yStartingValue = 45
)

func main() {
	defer fmt.Println("\nleaving main")
	col := readCSV("/Users/christianhovenbitzer/go/src/github.com/AnotherCoolDude/plotthatshit/heartbeat.csv")

	plot, err := plot.New()
	if err != nil {
		fmt.Println(err)
	}
	plot.Title.Text = "Auswertung"
	plot.X.Label.Text = "Zeit in Sekunden [s]"
	plot.Y.Label.Text = "Herzschläge pro Minute [bpm]"
	plot.Add(plotter.NewGrid())

	userMap := col.getUserMap()

	intmap := refactorMap(&userMap)
	for i, coords := range intmap {

		if len(userToPrint) != 0 && !contains(userToPrint, i) {
			continue
		}

		fmt.Printf("printing proband: %d\n", i)
		label := []string{}

		limitIndex := 0
		for i := range coords[0] {
			if coords[0][i] >= timeLimit {
				limitIndex = i
				break
			}
		}
		if timeLimit < 0 {
			limitIndex = len(coords[1]) - 1
		}
		xys := make(plotter.XYs, limitIndex)
		fmt.Println(limitIndex)
		for i := range coords[1][:limitIndex] {
			xys[i].X = coords[0][i]
			xys[i].Y = coords[1][i]

			label = append(label, fmt.Sprintf("%.1f\n%.1f ", coords[0][i], coords[1][i]))
		}
		label[0] = fmt.Sprintf("P%d\n", i) + label[0]
		xyLabeller := plotter.XYLabels{
			XYs:    xys,
			Labels: label,
		}
		l, s, err := plotter.NewLinePoints(xys)
		labelPlot, err := plotter.NewLabels(xyLabeller)
		if err != nil {
			fmt.Println(err)
		}
		rdColor := randomcolor.GetRandomColorInRgb()
		l.Color = color.RGBA{R: uint8(rdColor.Red), G: uint8(rdColor.Green), B: uint8(rdColor.Blue), A: 255}
		s.Color = color.RGBA{R: uint8(rdColor.Red), G: uint8(rdColor.Green), B: uint8(rdColor.Blue), A: 255}
		ts := make([]draw.TextStyle, xyLabeller.Len())
		for i := range labelPlot.TextStyle {
			ts[i] = smallDrawStyle
		}
		labelPlot.TextStyle = ts
		plot.Add(l, labelPlot)
		plot.Legend.Add(fmt.Sprintf("Proband: %d", i), l, s)
	}

	addStimuli(plot)
	plot.X.Tick.Marker = addTicks(plot.X.Tick.Marker)
	plot.Y.Tick.Marker = addTicks(plot.Y.Tick.Marker)
	setFrameFonts(plot)
	savePlot(plot)

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

type stimulus struct {
	xValue float64
	name   string
}

type legendPlacement struct {
	top     bool
	left    bool
	xOffset int
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

	}
	for _, u := range *col.userColl {
		fmt.Printf("userid: %s, dataamount: %d\n", u.id, len(*u.data))
	}

	return col

}

func savePlot(plot *plot.Plot) {
	wt, err := plot.WriterTo(imgWidth, imgHeight, "png")
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

type tickerFunc func(min, max float64) []plot.Tick

func (tkfn tickerFunc) Ticks(min, max float64) []plot.Tick { return tkfn(min, max) }

func addTicks(marker plot.Ticker) plot.Ticker {
	return tickerFunc(func(min, max float64) []plot.Tick {
		var out []plot.Tick
		interval := float64(1)
		for i := 0; i < int(math.Round(max)); i++ {
			nTick := plot.Tick{
				Value: interval * float64(i),
				Label: "",
			}
			if i%10 == 0 {
				nTick.Label = fmt.Sprintf("%.2f", interval*float64(i))
			}
			out = append(out, nTick)
		}
		return out
	})
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

func min(a, b int8) int8 {
	if a < b {
		return a
	}
	return b
}

func refactorMap(m *map[string][][]float64) map[int][][]float64 {
	intMap := map[int][][]float64{}
	re := regexp.MustCompile("[0-9]+")
	for k, v := range *m {
		fmt.Println(re.FindAllString(k, 1))
		i, err := strconv.Atoi(re.FindAllString(k, 1)[0])
		if err != nil {
			continue
		}
		intMap[i] = v
	}
	for i, k := range intMap {
		fmt.Printf("Key: %d, len(0): %d, len(1): %d\n", i, len(k[0]), len(k[1]))
	}
	return intMap
}

func addStimuli(plot *plot.Plot) {
	for i := 0; i < 9; i++ {

		if stimuli[i].xValue >= timeLimit {
			continue
		}

		xys := make(plotter.XYs, 2)
		xys[0].Y = yStartingValue
		xys[0].X = stimuli[i].xValue
		xys[1].Y = plot.Y.Max
		xys[1].X = stimuli[i].xValue
		l, _ := plotter.NewLine(xys)
		l.Dashes = []vg.Length{5, 5}
		l.Color = plotutil.Color(4)
		xyLabels := plotter.XYLabels{
			XYs:    xys,
			Labels: []string{stimuli[i].name, ""},
		}
		labels, _ := plotter.NewLabels(xyLabels)
		ts := make([]draw.TextStyle, xyLabels.Len())
		for i := range labels.TextStyle {
			ts[i] = largeDrawStyle
		}
		labels.TextStyle = ts
		plot.Add(l, labels)
	}

}

func contains(slice []int, number int) bool {
	for _, v := range slice {
		if v == number {
			return true
		}
	}
	return false
}

func setFrameFonts(plot *plot.Plot) {
	plot.Title.Font = lFont
	plot.X.Label.TextStyle = largeDrawStyle
	plot.Y.Label.TextStyle = largeDrawStyle
	plot.X.Tick.Label.Font = lFont
	plot.Y.Tick.Label.Font = lFont
	plot.Legend.Font = sFont
	plot.Legend.Top = lPlacement.top
	plot.Legend.Left = lPlacement.left
	plot.Legend.XOffs = vg.Length(lPlacement.xOffset)
}
