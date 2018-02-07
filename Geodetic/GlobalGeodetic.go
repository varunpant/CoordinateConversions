/*
TMS Global Geodetic Profile
---------------------------

Functions necessary for generation of global tiles in Plate Carre projection,
EPSG:4326, "unprojected profile".

Such tiles are compatible with Google Earth (as any other EPSG:4326 rasters)
and you can overlay the tiles on top of OpenLayers base map.

Pixel and tile coordinates are in TMS notation (origin [0,0] in bottom-left).

What coordinate conversions do we need for TMS Global Geodetic tiles?

Global Geodetic tiles are using geodetic coordinates (latitude,longitude)
directly as planar coordinates XY (it is also called Unprojected or Plate
Carre). We need only scaling to pixel pyramid and cutting to tiles.
Pyramid has on top level two tiles, so it is not square but rectangle.
Area [-180,-90,180,90] is scaled to 512x256 pixels.
TMS has coordinate origin (for pixels and tiles) in bottom-left corner.
Rasters are in EPSG:4326 and therefore are compatible with Google Earth.

LatLon      <->      Pixels      <->     Tiles

WGS84 coordinates   Pixels in pyramid  Tiles in pyramid
lat/lon         XY pixels Z zoom      XYZ from TMS
EPSG:4326
   .----.                ----
 /       \     <->    /--------/    <->      TMS
 \      /         /--------------/
   -----        /--------------------/
WMS, KML    Web Clients, Google Earth  TileMapService
*/
package Geodetic

import "math"

type GlobalGeodetic struct {
	tileSize float64
}

func NewGlobalGeodetic(tileSize int64) *GlobalGeodetic {
	return &GlobalGeodetic{
		tileSize: float64(tileSize),
	}
}

//Converts lat/lon to pixel coordinates in given zoom of the EPSG:4326 pyramid
func  (this *GlobalGeodetic)  LatLonToPixels(lat float64, lon float64, zoom int64) (float64, float64) {

	res := 180 / 256.0 / math.Pow(2, float64(zoom))
	px := (180 + lat) / res
	py := (90 + lon) / res
	return px, py

}

//Returns coordinates of the tile covering region in pixel coordinates
func (this *GlobalGeodetic) PixelsToTile(px float64, py float64) (int64, int64) {
	tx := int64(math.Ceil(px/this.tileSize) - 1)
	ty := int64(math.Ceil(py/this.tileSize) - 1)
	return tx, ty
}

func (this *GlobalGeodetic) Resolution(zoom int64) float64 {
	return 180 / 256.0 / math.Pow(2, float64(zoom))
}

//Returns bounds of the given tile
func (this *GlobalGeodetic) TileBounds(tx float64, ty float64, zoom int64) (float64, float64, float64, float64) {
	res := 180 / 256.0 / math.Pow(2, float64(zoom))
	return tx*256*res - 180, ty*256*res - 90, (tx+1)*256*res - 180, (ty+1)*256*res - 90
}
