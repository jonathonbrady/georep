package overpass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OverpassClient struct {
	Client *http.Client
	BoundingBoxes map[string]string
}

type OverpassResponse struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type  string   `json:"type"`
	ID    int64    `json:"id"`
	Lat   float64  `json:"lat,omitempty"`
	Lon   float64  `json:"lon,omitempty"`
	Nodes []int64  `json:"nodes,omitempty"`
}

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
		Client: http.DefaultClient,
		BoundingBoxes: boundingBoxes,
	}, nil
}

func (oc *OverpassClient) GetLocationsOnRoad(country string, road string) ([][2]float64, error) {
	bbox, ok := oc.BoundingBoxes[country]
	if !ok {
		return [][2]float64{}, fmt.Errorf("country %s not found in bounding boxes file", country)
	}

	query := fmt.Sprintf(`
	[out:json];
	way[highway]["ref"="%s"](%s);
	(._;>;);
	out body;
	`, road, bbox)

	resp, err := oc.Client.Post("https://overpass-api.de/api/interpreter", "application/x-www-form-urlencoded", strings.NewReader("data="+query))
	if err != nil {
		return [][2]float64{}, fmt.Errorf("failed to query Overpass API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var overpassResp OverpassResponse
	err = json.Unmarshal(body, &overpassResp)
	if err != nil {
		return [][2]float64{}, fmt.Errorf("failed to parse Overpass API response: %v", err)
	}

	coordinates := make([][2]float64, 0)
	for _, el := range overpassResp.Elements {
		if el.Type == "node" {
			coordinates = append(coordinates, [2]float64{el.Lat, el.Lon})
		}
	}

	if len(coordinates) == 0 {
		return [][2]float64{}, fmt.Errorf("no nodes found in Overpass API response")
	}
	return coordinates, nil
}
