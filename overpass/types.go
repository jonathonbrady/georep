package overpass

import "net/http"

type OverpassClient struct {
	Client        *http.Client
	BoundingBoxes map[string]string
}

type OverpassResponse struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type  string  `json:"type"`
	ID    int64   `json:"id"`
	Lat   float64 `json:"lat,omitempty"`
	Lon   float64 `json:"lon,omitempty"`
	Nodes []int64 `json:"nodes,omitempty"`
}

type Latlong struct {
	Latitude  float64
	Longitude float64
}
