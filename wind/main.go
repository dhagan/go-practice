package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// WindData stores time, direction, and velocity
type WindData struct {
	Timestamp time.Time
	Direction float64
	Velocity  float64
}

func main() {
	// Load wind data from CSV
	file, err := os.Open("wind_data_minute.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var data []WindData
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05", record[0])
		direction, _ := strconv.ParseFloat(record[1], 64)
		velocity, _ := strconv.ParseFloat(record[2], 64)
		data = append(data, WindData{timestamp, direction, velocity})
	}

	// Create plot
	p := plot.New()
	p.Title.Text = "Wind Direction and Velocity Over Time"
	p.X.Label.Text = "Value (Direction or Velocity)"
	p.Y.Label.Text = "Time"

	// Plot wind direction and velocity with time on the y-axis
	directionPoints := make(plotter.XYs, len(data))
	velocityPoints := make(plotter.XYs, len(data))
	for i, d := range data {
		directionPoints[i].X = d.Direction                 // Wind direction on x-axis
		directionPoints[i].Y = float64(d.Timestamp.Unix()) // Time on y-axis
		velocityPoints[i].X = d.Velocity                   // Wind velocity on x-axis
		velocityPoints[i].Y = float64(d.Timestamp.Unix())  // Time on y-axis
	}

	directionLine, err := plotter.NewLine(directionPoints)
	if err != nil {
		log.Fatal(err)
	}
	directionLine.Color = plotutil.Color(1)
	directionLine.LineStyle.Width = vg.Points(1)

	velocityLine, err := plotter.NewLine(velocityPoints)
	if err != nil {
		log.Fatal(err)
	}
	velocityLine.Color = plotutil.Color(2)
	velocityLine.LineStyle.Width = vg.Points(1)

	p.Add(directionLine, velocityLine)
	p.Legend.Add("Wind Direction", directionLine)
	p.Legend.Add("Wind Velocity", velocityLine)

	// Set Y-axis to display time formatted as HH:MM and invert it for waterfall effect
	p.Y.Tick.Marker = plot.TimeTicks{Format: "15:04"}
	//p.Y.Scale = plot.ReverseScale{}  // Invert y-axis

	// Save plot
	if err := p.Save(8*vg.Inch, 10*vg.Inch, "wind_waterfall_plot.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Plot saved as wind_waterfall_plot.png")
}
