package main

import (
    "github.com/tkrajina/gpxgo/gpx"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	
	gpxFileName := flag.String("gpx", "", "the gpx file (either specify --gpx or --latitude and --longitude)")
	latitude := flag.Float64("latitude", 0, "the latitude of the base point")
	longitude := flag.Float64("longitude", 0, "the longitude of the base point")
	flag.Bool("links", true, "Adds for each point a link to OSM with a marker (optional)")
	flag.Usage = func() {
		w := flag.CommandLine.Output() 
		fmt.Fprintf(w, "Usage of %s: (--gpx FILE || --latitude LAT --longitude LON) [--links] GPXFILES... \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "Output: Distance to FILE/LAT+LON in meters;"+
		"\n\tGFILE from GPXFILES;\n\tMIN LAT from GFILE;\n\tMIN LON from GFILE;\n\tTIME from GFILE;"+
		"\n\tMIN LAT from FILE or LAT;\n\tMIN LOT from FILE or LON;\n\tTIME from FILE or default\n\t[;link to GFILE;link to FILE/LAT+LON]\n")
	
	}
	
	flag.Parse()
	// see https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name]=true } )

	if !(flagset["gpx"] || (flagset["latitude"] && flagset["longitude"])) {
		log.Fatal("Either --gpx or both --latitude and --longitude must be given")
		return
	}
	var mainPoints []gpx.GPXPoint
	if(flagset["gpx"]){
		gpxFile, shouldReturn := parseFile(gpxFileName)
		if shouldReturn {
			return
		}
		mainPoints = getPointsFromFile(gpxFile)
	} else {
		basePoint := gpx.GPXPoint{Point: gpx.Point{Latitude: *latitude, Longitude: *longitude, Elevation: *gpx.NewNullableFloat64(0)}}
		mainPoints = append(mainPoints, basePoint)
	}
	
	//printPoints(mainPoints)
	for _, filename := range flag.Args() {
		var filePoints [] gpx.GPXPoint
		newGpxFile, shouldReturn2 := parseFile(&filename)
		if shouldReturn2 {
			return
		}
		filePoints = getPointsFromFile(newGpxFile)
		var min float64
		var minLat1, minLon1, minLat2, minLon2 float64
		var time1, time2 time.Time
		min = 100000000.
		for _, point := range filePoints {
			for _, point2 := range mainPoints {
				var dist = point.Distance2D(&point2)
				if(dist < min){
					min = dist
					minLat1 = point.Latitude
					minLon1 = point.Longitude
					time1 = point.Timestamp
					minLat2 = point2.Latitude
					minLon2 = point2.Longitude
					time2 = point2.Timestamp
				}
			}				
		}
		// 
		var links string
		if(flagset["links"]){
			var link string 
			var link2 string
			link = fmt.Sprintf("http://www.openstreetmap.org/?mlat=%f&mlon=%f&layers=M", minLat1, minLon1)
			link2 = fmt.Sprintf("http://www.openstreetmap.org/?mlat=%f&mlon=%f&layers=M", minLat2, minLon2)
			links = ";" + link + ";" + link2
		} else {
			links = ""
		}
		fmt.Printf("%.2f;%s;%f;%f;%s;%f;%f;%s%s\n", min, filename, minLat1,minLon1, time1.String(), minLat2, minLon2,time2.String(), links)
	}
}

func parseFile(fileName *string) (*gpx.GPX, bool) {
	file, err := os.Open(*fileName)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)

	gpxFile, err := gpx.ParseBytes(byteValue)
	if err != nil {
		log.Println("Err: " + string(err.Error()))
		return nil, true
	}
	return gpxFile, false
}

func printPoints(mainPoints []*gpx.GPXPoint) {
	for _, point := range mainPoints {
		fmt.Print(*point)
		fmt.Print("\n")
	}
}

func getPointsFromFile(gpxFile *gpx.GPX) ([]gpx.GPXPoint) {
	var points []gpx.GPXPoint
	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				points = append(points, point)
			}
		}
	}
	return points
}

