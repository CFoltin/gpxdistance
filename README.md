# Goal

The tool `gpxdistance` is able to calculate the minimum distance of a list of GPX tracks to either a single point or to a another track. 

The use cases are 
- test if one of the GPX cases are close to a given point (Question: "I have been there already...")
- find common GPX tracks (e.g. to find a tracked file with a set of commercial tracks)

# Build
```
go test
go build gpxdistance.go 
```

# Usage
```
Usage of gpxdistance: (--gpx FILE || --latitude LAT --longitude LON) [--links] GPXFILES... 
  -gpx string
        the gpx file (either specify --gpx or --latitude and --longitude)
  -latitude float
        the latitude of the base point
  -links
        Adds for each point a link to OSM with a marker (optional) (default true)
  -longitude float
        the longitude of the base point
Output: Distance to FILE/LAT+LON in meters;
        GFILE from GPXFILES;
        MIN LAT from GFILE;
        MIN LON from GFILE;
        TIME from GFILE;
        MIN LAT from FILE or LAT;
        MIN LOT from FILE or LON;
        TIME from FILE or default
        [;link to GFILE;link to FILE/LAT+LON]
```

