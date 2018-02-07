/*
    TMS Global Mercator Profile
	---------------------------

	Functions necessary for generation of tiles in Spherical Mercator projection,
	EPSG:900913 (EPSG:gOOglE, Google Maps Global Mercator), EPSG:3785, OSGEO:41001.

	Such tiles are compatible with Google Maps, Microsoft Virtual Earth, Yahoo Maps,
	UK Ordnance Survey OpenSpace API, ...
	and you can overlay them on top of base maps of those web mapping applications.

	Pixel and tile coordinates are in TMS notation (origin [0,0] in bottom-left).

	What coordinate conversions do we need for TMS Global Mercator tiles::

	     LatLon      <->       Meters      <->     Pixels    <->       Tile

	 WGS84 coordinates   Spherical Mercator  Pixels in pyramid  Tiles in pyramid
	     lat/lon            XY in metres     XY pixels Z zoom      XYZ from TMS
	    EPSG:4326           EPSG:900913
	     .----.              ---------               --                TMS
	    /      \     <->     |       |     <->     /----/    <->      Google
	    \      /             |       |           /--------/          QuadTree
	     -----               ---------         /------------/
	   KML, public         WebMapService         Web Clients      TileMapService

	What is the coordinate extent of Earth in EPSG:900913?

	  [-20037508.342789244, -20037508.342789244, 20037508.342789244, 20037508.342789244]
	  Constant 20037508.342789244 comes from the circumference of the Earth in meters,
	  which is 40 thousand kilometers, the coordinate origin is in the middle of extent.
      In fact you can calculate the constant as: 2 * math.pi * 6378137 / 2.0
	  $ echo 180 85 | gdaltransform -s_srs EPSG:4326 -t_srs EPSG:900913
	  Polar areas with abs(latitude) bigger then 85.05112878 are clipped off.

	What are zoom level constants (pixels/meter) for pyramid with EPSG:900913?

	  whole region is on top of pyramid (zoom=0) covered by 256x256 pixels tile,
	  every lower zoom level resolution is always divided by two
	  initialResolution = 20037508.342789244 * 2 / 256 = 156543.03392804062

	What is the difference between TMS and Google Maps/QuadTree tile name convention?

	  The tile raster itself is the same (equal extent, projection, pixel size),
	  there is just different identification of the same raster tile.
	  Tiles in TMS are counted from [0,0] in the bottom-left corner, id is XYZ.
	  Google placed the origin [0,0] to the top-left corner, reference is XYZ.
	  Microsoft is referencing tiles by a QuadTree name, defined on the website:
	  http://msdn2.microsoft.com/en-us/library/bb259689.aspx

	The lat/lon coordinates are using WGS84 datum, yeh?

	  Yes, all lat/lon we are mentioning should use WGS84 Geodetic Datum.
	  Well, the web clients like Google Maps are projecting those coordinates by
	  Spherical Mercator, so in fact lat/lon coordinates on sphere are treated as if
	  the were on the WGS84 ellipsoid.

	  From MSDN documentation:
	  To simplify the calculations, we use the spherical form of projection, not
	  the ellipsoidal form. Since the projection is used only for map display,
	  and not for displaying numeric coordinates, we don't need the extra precision
	  of an ellipsoidal projection. The spherical projection causes approximately
	  0.33 percent scale distortion in the Y direction, which is not visually noticable.

	How do I create a raster in EPSG:900913 and convert coordinates with PROJ.4?

	  You can use standard GIS tools like gdalwarp, cs2cs or gdaltransform.
	  All of the tools supports -t_srs 'epsg:900913'.

	  For other GIS programs check the exact definition of the projection:
	  More info at http://spatialreference.org/ref/user/google-projection/
	  The same projection is degined as EPSG:3785. WKT definition is in the official
	  EPSG database.

	  Proj4 Text:
	    +proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0
	    +k=1.0 +units=m +nadgrids=@null +no_defs

	  Human readable WKT format of EPGS:900913:
	     PROJCS["Google Maps Global Mercator",
	         GEOGCS["WGS 84",
	             DATUM["WGS_1984",
	                 SPHEROID["WGS 84",6378137,298.2572235630016,
	                     AUTHORITY["EPSG","7030"]],
	                 AUTHORITY["EPSG","6326"]],
	             PRIMEM["Greenwich",0],
	             UNIT["degree",0.0174532925199433],
	             AUTHORITY["EPSG","4326"]],
	         PROJECTION["Mercator_1SP"],
	         PARAMETER["central_meridian",0],
	         PARAMETER["scale_factor",1],
	         PARAMETER["false_easting",0],
	         PARAMETER["false_northing",0],
	         UNIT["metre",1,
	             AUTHORITY["EPSG","9001"]]]
*/
package Mercator

import "math"

type GlobalMercator struct {
	tileSize          float64
	initialResolution float64
	originShift       float64
}

func NewGlobalMercator(tileSize float64) *GlobalMercator {

	return &GlobalMercator{
		tileSize,
		2 * math.Pi * 6378137 / tileSize,
		//156543.03392804062 for tileSize 256 pixels
		2 * math.Pi * 6378137 / 2.0,
	}
}

//Converts given lat/lon in WGS84 Datum to XY in Spherical Mercator EPSG:900913
func (this *GlobalMercator) LatLonToMeters(lat float64, lon float64) (float64, float64) {
	mx := lon * this.originShift / 180.0
	my := math.Log(math.Tan((90+lat)*math.Pi/360.0)) / (math.Pi / 180.0)

	my = my * this.originShift / 180.0
	return mx, my
}

//Converts XY point from Spherical Mercator EPSG:900913 to lat/lon in WGS84 Datum
func (this *GlobalMercator) MetersToLatLon(mx float64, my float64) (float64, float64) {

	lon := (mx / this.originShift) * 180.0
	lat := (my / this.originShift) * 180.0

	lat = 180 / math.Pi * (2*math.Atan(math.Exp(lat*math.Pi/180.0)) - math.Pi/2.0)
	return lat, lon
}

//Converts pixel coordinates in given zoom level of pyramid to EPSG:900913
func (this *GlobalMercator) PixelsToMeters(px float64, py float64, zoom int64) (float64, float64) {

	res := this.Resolution(zoom)
	mx := px*res - this.originShift
	my := py*res - this.originShift
	return mx, my

}

//Converts EPSG:900913 to pyramid pixel coordinates in given zoom level
func (this *GlobalMercator) MetersToPixels(mx float64, my float64, zoom int64) (float64, float64) {

	res := this.Resolution(zoom)
	px := (mx + this.originShift) / res
	py := (my + this.originShift) / res
	return px, py
}

//Returns a tile covering region in given pixel coordinates
func (this *GlobalMercator) PixelsToTile(px float64, py float64) (int64, int64) {
	tx := int64(math.Ceil(px/this.tileSize) - 1)
	ty := int64(math.Ceil(py/this.tileSize) - 1)
	return tx, ty
}

//Move the origin of pixel coordinates to top-left corner
func (this *GlobalMercator) PixelsToRaster(px float64, py float64, zoom int64) (float64, float64) {
	mapSize := uint(this.tileSize) << uint(zoom)
	return px, float64(mapSize) - py
}

//Returns tile for given mercator coordinates
func (this *GlobalMercator) MetersToTile(mx float64, my float64, zoom int64) (int64, int64) {
	px, py := this.MetersToPixels(mx, my, zoom)
	return this.PixelsToTile(px, py)
}

//Returns bounds of the given tile in EPSG:900913 coordinates
func (this *GlobalMercator) TileBounds(tx float64, ty float64, zoom int64) (float64, float64, float64, float64) {
	minx, miny := this.PixelsToMeters(tx*this.tileSize, ty*this.tileSize, zoom)
	maxx, maxy := this.PixelsToMeters((tx+1)*this.tileSize, (ty+1)*this.tileSize, zoom)
	return minx, miny, maxx, maxy
}

//Returns bounds of the given tile in latutude/longitude using WGS84 datum
func (this *GlobalMercator) TileLatLonBounds(tx float64, ty float64, zoom int64) (float64, float64, float64, float64) {

	minx, miny, maxx, maxy := this.TileBounds(tx, ty, zoom)
	minLat, minLon := this.MetersToLatLon(minx, miny)
	maxLat, maxLon := this.MetersToLatLon(maxx, maxy)

	return minLat, minLon, maxLat, maxLon

}

//Resolution (meters/pixel) for given zoom level (measured at Equator)
func (this *GlobalMercator) Resolution(zoom int64) float64 {
	//return (2 * math.pi * 6378137) / (self.tileSize * 2**zoom)
	return this.initialResolution / math.Pow(2, float64(zoom))
}

//Maximal scaledown zoom of the pyramid closest to the pixelSize.
func (this *GlobalMercator) ZoomForPixelSize(pixelSize float64 ) int{
	for i :=0;i<30;i++{
		if pixelSize > this.Resolution(i){
			if i!=0{
				return 0
			}else{
				return i-1
			}
		}
	}
}


