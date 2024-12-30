package googlemaps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func NewGoogleMapsClient() (*GoogleMapsClient, error) {
	key, ok := os.LookupEnv("GOOGLE_MAPS_API_KEY")
	if !ok {
		return nil, fmt.Errorf("google maps api key environment variable not set")
	}

	return &GoogleMapsClient{
		Client:   http.DefaultClient,
		Auth:     key,
		APICalls: make(map[string]int),
	}, nil
}

// Snap locations to the nearest road. A maximum of 100 locations will be used.
func (gc *GoogleMapsClient) SnapToRoads(locations [][2]float64) ([][2]float64, error) {
	if calls, ok := gc.APICalls["SnapToRoads"]; ok {
		gc.APICalls["SnapToRoads"] = calls + 1
	} else {
		gc.APICalls["SnapToRoads"] = 1
	}

	strs := make([]string, 0)
	for _, location := range locations {
		strs = append(strs, fmt.Sprintf("%f%%2C%f", location[0], location[1]))
	}
	path := strings.Join(strs, "%7C")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://roads.googleapis.com/v1/nearestRoads?points=%s&key=%s", path, gc.Auth), http.NoBody)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.Client.Do(req)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return [][2]float64{}, fmt.Errorf("bad status from snap to roads API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("reading response body: %v", err)
	}

	var response SnapToRoadsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("unmarshaling response: %v", err)
	}

	if len(response.SnappedPoints) == 0 {
		return [][2]float64{{0, 0}}, nil
	}

	snapped := make([][2]float64, 0)
	for _, point := range response.SnappedPoints {
		p := [2]float64{
			point.Location.Latitude,
			point.Location.Longitude,
		}
		snapped = append(snapped, p)
	}

	return snapped, nil
}

// Locations should not be selected where there is no official Google Street View coverage.
func (gc *GoogleMapsClient) ValidateCoverage(latlong [2]float64) (bool, error) {
	if calls, ok := gc.APICalls["Metadata"]; ok {
		gc.APICalls["Metadata"] = calls + 1
	} else {
		gc.APICalls["Metadata"] = 1
	}

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
