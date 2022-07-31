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
