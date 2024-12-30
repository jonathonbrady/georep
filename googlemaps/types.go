package googlemaps

import "net/http"

type GoogleMapsClient struct {
	Client *http.Client
	Auth   string

	APICalls map[string]int
}

type GetMetadataResponse struct {
	Copyright string `json:"copyright"`
	Date      string `json:"date"`
	Location  struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"lng"`
	} `json:"location"`
	PanoId string `json:"pano_id"`
	Status string `json:"status"`
}

type SnapToRoadsResponse struct {
	SnappedPoints []struct {
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		OriginalIndex int    `json:"originalIndex,omitempty"`
		PlaceID       string `json:"placeId"`
	} `json:"snappedPoints"`
}
