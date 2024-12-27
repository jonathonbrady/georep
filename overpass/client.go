package overpass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func NewOverpassClient() (*OverpassClient, error) {
	file, err := os.ReadFile("D:/geosr/overpass/bounding_boxes.json")
	if err != nil {
		return nil, fmt.Errorf("reading bounding boxes file: %v", err)
	}

	var boundingBoxes map[string]string
	err = json.Unmarshal(file, &boundingBoxes)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling bounding boxes file: %v", err)
	}

	return &OverpassClient{
		Client:        http.DefaultClient,
		BoundingBoxes: boundingBoxes,
	}, nil
}

func (oc *OverpassClient) GetLocationsOnRoad(country string, road string) ([]Latlong, error) {
	bbox, ok := oc.BoundingBoxes[country]
	if !ok {
		return []Latlong{}, fmt.Errorf("country %s not found in bounding boxes file", country)
	}

	query := fmt.Sprintf(`
	[out:json];
	way[highway]["ref"="%s"](%s);
	(._;>;);
	out body;
	`, road, bbox)

	resp, err := oc.Client.Post("https://overpass-api.de/api/interpreter", "application/x-www-form-urlencoded", strings.NewReader("data="+query))
	if err != nil {
		return []Latlong{}, fmt.Errorf("failed to query Overpass API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Latlong{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var overpassResp OverpassResponse
	err = json.Unmarshal(body, &overpassResp)
	if err != nil {
		return []Latlong{}, fmt.Errorf("failed to parse Overpass API response: %v", err)
	}

	coordinates := make([]Latlong, 0)
	for _, el := range overpassResp.Elements {
		if el.Type == "node" {
			coordinates = append(coordinates, Latlong{el.Lat, el.Lon})
		}
	}

	if len(coordinates) == 0 {
		return []Latlong{}, fmt.Errorf("no nodes found in Overpass API response")
	}
	return coordinates, nil
}
