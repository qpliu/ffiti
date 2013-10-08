package data

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const (
	dlat = 0.0005
	dlng = 0.0002
)

type Location struct {
	Latitude         float64  `json:"lat"`
	Longitude        float64  `json:"lng"`
	Altitude         *float64 `json:"alt,omitempty"`
	Accuracy         *float64 `json:"acc,omitempty"`
	AltitudeAccuracy *float64 `json:"altacc,omitempty"`
	Alpha            *float64 `json:"alpha,omitempty"`
	Beta             *float64 `json:"beta,omitempty"`
	Gamma            *float64 `json:"gamma,omitempty"`
}

func parseParam(r *http.Request, param string) (*float64, error) {
	value, err := strconv.ParseFloat(r.FormValue(param), 64)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func GetLocation(r *http.Request) (*Location, error) {
	lat, err := parseParam(r, "lat")
	if err != nil {
		return nil, err
	}
	lng, err := parseParam(r, "lng")
	if err != nil {
		return nil, err
	}
	alt, _ := parseParam(r, "alt")
	acc, _ := parseParam(r, "acc")
	altacc, _ := parseParam(r, "altacc")
	alpha, _ := parseParam(r, "alpha")
	beta, _ := parseParam(r, "beta")
	gamma, _ := parseParam(r, "gamma")
	return &Location{
		Latitude:         *lat,
		Longitude:        *lng,
		Altitude:         alt,
		Accuracy:         acc,
		AltitudeAccuracy: altacc,
		Alpha:            alpha,
		Beta:             beta,
		Gamma:            gamma,
	}, nil
}

func (loc Location) key() (int, int) {
	return int(loc.Latitude / dlat), int(loc.Longitude / dlng)
}

func (loc Location) Key() string {
	x, y := loc.key()
	return fmt.Sprintf("%v,%v", x, y)
}

func (loc Location) Keys() [9]string {
	x, y := loc.key()
	return [9]string{
		fmt.Sprintf("%v,%v", x, y),
		fmt.Sprintf("%v,%v", x+1, y),
		fmt.Sprintf("%v,%v", x-1, y),
		fmt.Sprintf("%v,%v", x, y+1),
		fmt.Sprintf("%v,%v", x+1, y+1),
		fmt.Sprintf("%v,%v", x-1, y+1),
		fmt.Sprintf("%v,%v", x, y-1),
		fmt.Sprintf("%v,%v", x+1, y-1),
		fmt.Sprintf("%v,%v", x-1, y-1),
	}
}

func (loc Location) Bounds() [4]float64 {
	lat := math.Floor(loc.Latitude/dlat) * dlat
	lng := math.Floor(loc.Longitude/dlng) * dlng
	return [4]float64{lat - 0.5*dlat, lng - 0.5*dlng, lat + 1.5*dlat, lng + 1.5*dlng}
}
