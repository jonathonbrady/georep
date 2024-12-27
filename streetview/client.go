package streetview

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func NewStreetViewClient() (*StreetViewClient, error) {
	key, ok := os.LookupEnv("GOOGLE_MAPS_API_KEY")
	if !ok {
		return nil, fmt.Errorf("google maps api key environment variable not set")
	}

	return &StreetViewClient{
		Client: http.DefaultClient,
		Auth:   key,
	}, nil
}

// Locations should not be selected where there is no official Google Street View coverage.
func (gc *StreetViewClient) ValidateCoverage(latlong [2]float64) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://maps.googleapis.com/maps/api/streetview/metadata?location=%f,%%20%f&key=%s", latlong[0], latlong[1], gc.Auth), http.NoBody)
	if err != nil {
		return false, fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad status from challenges API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body: %v", err)
	}

	var response GetMetadataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("unmarshaling response: %v", err)
	}

	// Will return ZERO_RESULTS if there is no coverage.
	if response.Status != "OK" {
		return false, nil
	}

	// Third-party coverage will not be copyright by Google.
	if response.Copyright != "© Google" {
		return false, nil
	}

	return true, nil
}
