package Geodetic

import (
	"testing"
	"math"

)

func floatEquals(a, b float64, eps float64) bool {
	if math.Abs(a-b) < eps {
		return true
	}
	return false
}

func TestLatLonToPixels(t *testing.T) {
	px, py := LatLonToPixels(51.5287718, -0.2416819, 2)

	if !floatEquals(px, 1317.141, 0.001) {
		t.Errorf("expected 1317.141 as px")
	}
	if !floatEquals(py, 510.625, 0.001) {
		t.Errorf("expected 510.625 as py ")
	}
}

func TestPixelsToTile(t *testing.T) {
	tx, ty := PixelsToTile(1317.141, 510.625, 256)

	if tx != 5 {
		t.Errorf("expected 5 as tx")
	}
	if ty != 1 {
		t.Errorf("expected 1 as ty")
	}
}

func TestResolution(t *testing.T) {
	expected_resolutions := [21]float64{
		0.703125,
		0.3515625,
		0.17578125,
		0.087890625,
		0.0439453125,
		0.02197265625,
		0.010986328125,
		0.0054931640625,
		0.00274658203125,
		0.001373291015625,
		0.0006866455078125,
		0.00034332275390625,
		0.000171661376953125,
		8.58306884765625e-05,
		4.291534423828125e-05,
		2.1457672119140625e-05,
		1.0728836059570312e-05,
		5.364418029785156e-06,
		2.682209014892578e-06,
		1.341104507446289e-06,
		6.705522537231445e-07,
	}

	for i := 0; i < 21; i++ {

		resolution := Resolution(int64(i))
		if !floatEquals(resolution, expected_resolutions[i],  0.0001) {
			t.Errorf("expected resolution %g at zoom %d but was %g", expected_resolutions[i], i, resolution)
		}
	}
}
