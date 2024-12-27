package streetview

import "net/http"

type StreetViewClient struct {
	Client *http.Client
	Auth   string
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
