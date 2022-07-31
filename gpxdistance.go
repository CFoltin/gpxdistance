package main

import (
    "github.com/tkrajina/gpxgo/gpx"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os"
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
	var mainPoints []*gpx.GPXPoint
	if(flagset["gpx"]){
		gpxFile, shouldReturn := parseFile(gpxFileName)
		if shouldReturn {
			return
		}
		mainPoints = getPointsFromFile(gpxFile)
	} else {
		basePoint := gpx.GPXPoint{Point: gpx.Point{Latitude: *latitude, Longitude: *longitude, Elevation: *gpx.NewNullableFloat64(0)}}
		mainPoints = append(mainPoints, &basePoint)
	}
	
	//printPoints(mainPoints)
	for _, filename := range flag.Args() {
		var filePoints [] *gpx.GPXPoint
		newGpxFile, shouldReturn2 := parseFile(&filename)
		if shouldReturn2 {
			return
		}
		filePoints = getPointsFromFile(newGpxFile)
		var min float64
		var minPoint *gpx.GPXPoint
		var minPoint2 *gpx.GPXPoint
		min = 100000000.
		minPoint = nil
		minPoint2 = nil
		for _, point := range filePoints {
			for _, point2 := range mainPoints {
				var dist = point.Distance2D(point2)
				if(dist < min){
					min = dist
					minPoint = point
					minPoint2 = point2
				}
			}				
		}
		// 
		var links string
		if(flagset["links"]){
			var link string 
			var link2 string
			link = fmt.Sprintf("http://www.openstreetmap.org/?mlat=%f&mlon=%f&layers=M", minPoint.Latitude, minPoint.Longitude)
			link2 = fmt.Sprintf("http://www.openstreetmap.org/?mlat=%f&mlon=%f&layers=M", minPoint2.Latitude, minPoint2.Longitude)
			links = ";" + link + ";" + link2
		} else {
			links = ""
		}
		fmt.Printf("%.2f;%s;%f;%f;%s;%f;%f;%s%s\n", min, filename, minPoint.Latitude,minPoint.Longitude,minPoint.Timestamp.String(), minPoint2.Latitude,minPoint2.Longitude,minPoint2.Timestamp.String(), links)
	}
}

func parseFile(jsonFileName *string) (*gpx.GPX, bool) {
	jsonFile, err := os.Open(*jsonFileName)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

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

func getPointsFromFile(gpxFile *gpx.GPX) ([]*gpx.GPXPoint) {
	var points []*gpx.GPXPoint
	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				points = append(points, &point)
			}
		}
	}
	return points
}

