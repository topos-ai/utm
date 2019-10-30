package utm

import (
	"math"
)

const (

	// UTM constants
	utmFalseEasting                    float64 = 500000
	utmSouthernHemisphereFalseNorthing float64 = 10000000

	// WGS84 constants
	wgs84SemiMajorAxis     float64 = 6378137
	wgs84InverseFlattening float64 = 298.257223563
	wgs84ScaleFactor       float64 = 0.9996

	// Semi-major axis
	a float64 = wgs84SemiMajorAxis

	// Semi-minor axis
	b float64 = (a / (1 + n)) * (1 + n2/4 + n4/64)

	// Flattening
	f float64 = 1 / wgs84InverseFlattening

	// False easting
	fe float64 = utmFalseEasting

	// Scale factor at natural origin
	kO float64 = wgs84ScaleFactor

	n  float64 = f / (2 - f)
	n2 float64 = n * n
	n3 float64 = n * n * n
	n4 float64 = n * n * n * n

	h1 float64 = n/2 - n2*2/3 + n3*5/16 + n4*41/180
	h2 float64 = n2*13/48 - n3*3/5 + n4*557/1440
	h3 float64 = n3*61/240 - n4*103/140
	h4 float64 = n4 * 49561 / 161280

	h1Prime float64 = n/2 - n2*2/3 + n3*37/96 - n4/360
	h2Prime float64 = n2/48 + n3/15 - n4*437/1440
	h3Prime float64 = n3*17/480 - n4*37/840
	h4Prime float64 = n4 * 4397 / 161280
)

var (

	// Excentricity
	e float64
)

func init() {
	e = math.Sqrt(2*f - f*f)
}

// λO returns the longitude in radians of the natural origin for the given UTM
// zone.
func λO(zone int) float64 {
	return float64(zone)*math.Pi/30 - math.Pi*61/60
}

// Radians returns the latitude and longitude values for the given projected
// coordinates in radians.
func Radians(zone int, south bool, easting, northing float64) (float64, float64) {

	// False northing
	fn := 0.
	if south {
		fn = utmSouthernHemisphereFalseNorthing
	}

	ξPrime := (northing - fn) / (b * kO)
	ηPrime := (easting - fe) / (b * kO)

	ξ1Prime := h1Prime * math.Sin(2*ξPrime) * math.Cosh(2*ηPrime)
	ξ2Prime := h2Prime * math.Sin(4*ξPrime) * math.Cosh(4*ηPrime)
	ξ3Prime := h3Prime * math.Sin(6*ξPrime) * math.Cosh(6*ηPrime)
	ξ4Prime := h4Prime * math.Sin(8*ξPrime) * math.Cosh(8*ηPrime)
	ξ0Prime := ξPrime - (ξ1Prime + ξ2Prime + ξ3Prime + ξ4Prime)

	η1Prime := h1Prime * math.Cos(2*ξPrime) * math.Sinh(2*ηPrime)
	η2Prime := h2Prime * math.Cos(4*ξPrime) * math.Sinh(4*ηPrime)
	η3Prime := h3Prime * math.Cos(6*ξPrime) * math.Sinh(6*ηPrime)
	η4Prime := h4Prime * math.Cos(8*ξPrime) * math.Sinh(8*ηPrime)
	η0Prime := ηPrime - (η1Prime + η2Prime + η3Prime + η4Prime)

	βPrime := math.Asin(math.Sin(ξ0Prime) / math.Cosh(η0Prime))
	qPrime := math.Asinh(math.Tan(βPrime))

	qDoublePrime := qPrime + e*math.Atanh(e*math.Tanh(qPrime))

	// Iterate until the difference is insignificant.
	for {
		nextqDoublePrime := qPrime + e*math.Atanh(e*math.Tanh(qDoublePrime))
		if nextqDoublePrime == qDoublePrime {
			break
		}

		qDoublePrime = nextqDoublePrime
	}

	ϕ := math.Atan(math.Sinh(qDoublePrime))
	λ := λO(zone) + math.Asin(math.Tanh(η0Prime)/math.Cos(βPrime))
	return ϕ, λ
}

// Project returns the projected easting and northing values for the given
// latitude and longitude.
func Project(zone int, south bool, latitude, longitude float64) (float64, float64) {

	// False northing
	fn := 0.
	if south {
		fn = utmSouthernHemisphereFalseNorthing
	}

	h1 := n/2 - n2*2/3 + n3*5/16 + n4*41/180
	h2 := n2*13/48 - n3*3/5 + n4*557/1440
	h3 := n3*61/240 - n4*103/140
	h4 := n4 * 49561 / 161280

	q := math.Asinh(math.Tan(latitude)) - e*math.Atanh(e*math.Sin(latitude))
	β := math.Atan(math.Sinh(q))

	η0 := math.Atanh(math.Cos(β) * math.Sin(longitude-λO(zone)))
	ξ0 := math.Asin(math.Sin(β) * math.Cosh(η0))

	η1 := h1 * math.Cos(2*ξ0) * math.Sinh(2*η0)
	η2 := h2 * math.Cos(4*ξ0) * math.Sinh(4*η0)
	η3 := h3 * math.Cos(6*ξ0) * math.Sinh(6*η0)
	η4 := h4 * math.Cos(8*ξ0) * math.Sinh(8*η0)

	ξ1 := h1 * math.Sin(2*ξ0) * math.Cosh(2*η0)
	ξ2 := h2 * math.Sin(4*ξ0) * math.Cosh(4*η0)
	ξ3 := h3 * math.Sin(6*ξ0) * math.Cosh(6*η0)
	ξ4 := h4 * math.Sin(8*ξ0) * math.Cosh(8*η0)

	ξ := ξ0 + ξ1 + ξ2 + ξ3 + ξ4
	η := η0 + η1 + η2 + η3 + η4

	e := fe + kO*b*η
	n := fn + kO*b*ξ
	return e, n
}
