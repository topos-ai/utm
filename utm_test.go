package utm

import (
	"math"
	"testing"
)

const (
	testDegreePrecision     float64 = .000000001
	testProjectionPrecision float64 = .000001
)

func TestRadians(t *testing.T) {
	tests := []struct {
		zone      int
		south     bool
		easting   float64
		northing  float64
		latitude  float64
		longitude float64
	}{
		// echo 580741 4504692 | cs2cs +proj=utm +zone=18 +ellps=WGS84 +to +proj=latlong -d 9
		{18, false, 580741, 4504692, 40.689170987, -74.044439711}, // Statue of Liberty
		// echo 206895 3009276 | cs2cs +proj=utm +zone=44 +ellps=WGS84 +to +proj=latlong -d 9
		{44, false, 206895, 3009276, 27.175000872, 78.041942295}, // Taj Mahal
		// echo 334799 6252086 | cs2cs +proj=utm +zone=56 +south +ellps=WGS84 +to +proj=latlong -d 9
		{56, true, 334799, 6252086, -33.858611864, 151.214164458}, // Sydney Opera House
		// echo 683471 7460682 | cs2cs +proj=utm +zone=23 +south +ellps=WGS84 +to +proj=latlong -d 9
		{23, true, 683471, 7460682, -22.951948534, -43.210553081}, // Christ the Redeemer
		// echo 539203 1358223 | cs2cs +proj=utm +zone=58 +south +ellps=WGS84 +to +proj=latlong -d 9
		{58, true, 539203, 1358223, -77.846321004, 166.668248472}, // McMurdo Station
	}

	for _, test := range tests {
		latitude, longitude := Radians(test.zone, test.south, test.easting, test.northing)
		easting, northing := Project(test.zone, test.south, latitude, longitude)

		latitude *= 180 / math.Pi
		if latitude-test.latitude >= testDegreePrecision {
			t.Errorf("expected latitude %f, got %f", test.latitude, latitude)
		}

		longitude *= 180 / math.Pi
		if longitude-test.longitude >= testDegreePrecision {
			t.Errorf("expected longitude %f, got %f", test.longitude, longitude)
		}

		if easting-test.easting >= testProjectionPrecision {
			t.Errorf("expected easting %f, got %f", test.easting, easting)
		}

		if easting-test.northing >= testProjectionPrecision {
			t.Errorf("expected northing %f, got %f", test.northing, northing)
		}
	}
}

func BenchmarkRadians(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Radians(18, false, 580741, 4504692)
	}
}

func BenchmarkProject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Radians(18, false, 40.689170987, -74.044439711)
	}
}
