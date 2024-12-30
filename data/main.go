package data

import (
	"encoding/json"
	"fmt"
	"georep/googlemaps"
	"math"
	"math/rand/v2"
	"os"
)

var NULL_LOCATION = [2]float64{0, 0}
var NO_LOCATIONS = [][2]float64{{0, 0}}

func generateRandomLocationsInSubdivision(country string, subdivision string) ([][2]float64, error) {
	file, err := os.ReadFile("D:/geosr/data/geojson/ne_10m_admin_1_states_provinces.json")
	if err != nil {
		return NO_LOCATIONS, err
	}

	var subdivisions GeoJSON
	err = json.Unmarshal(file, &subdivisions)
	if err != nil {
		return NO_LOCATIONS, err
	}

	for _, feature := range subdivisions.Features {
		if country == feature.Properties.Admin && subdivision == feature.Properties.NameEn {
			if feature.Geometry.Type == "Polygon" {
				var coordinates [][][2]float64
				if err = json.Unmarshal(feature.Geometry.Coordinates, &coordinates); err != nil {
					return NO_LOCATIONS, err
				}

				// The GeoJSON is (long, lat) instead of (lat, long). Why?
				actual := make([][2]float64, 0)
				for _, coordinate := range coordinates[0] {
					actual = append(actual, [2]float64{coordinate[1], coordinate[0]})
				}

				locations := make([][2]float64, 0)
				for i := 0; i < 100; i++ {
					location := getRandomPointInPolygon(actual)
					locations = append(locations, location)
				}
				return locations, nil
			}
		}
	}

	return NO_LOCATIONS, fmt.Errorf("subdivision %v likely does not exist in %v", subdivision, country)
}

func isPointInPolygon(point [2]float64, polygon [][2]float64) bool {
	n := len(polygon)
	inside := false
	x, y := point[0], point[1]

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := polygon[i][0], polygon[i][1]
		xj, yj := polygon[j][0], polygon[j][1]

		intersect := ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
	}
	return inside
}

func getBoundingBox(polygon [][2]float64) ([2]float64, [2]float64) {
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	for _, p := range polygon {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	return [2]float64{minX, minY}, [2]float64{maxX, maxY}
}

func getRandomPointInPolygon(polygon [][2]float64) [2]float64 {
	min, max := getBoundingBox(polygon)

	for {
		x := min[0] + rand.Float64()*(max[0]-min[0])
		y := min[1] + rand.Float64()*(max[1]-min[1])
		point := [2]float64{x, y}

		if isPointInPolygon(point, polygon) {
			return point
		}
	}
}

func GetLocationsInSubdivision(country string, subdivision string, count int, sv *googlemaps.GoogleMapsClient) ([][2]float64, error) {
	locations := make([][2]float64, 0)

	for len(locations) < count {
		// Generate 100 locations within the polygon defined by the GeoJSON for this subdivision.
		randomLocations, err := generateRandomLocationsInSubdivision(country, subdivision)
		if err != nil {
			return [][2]float64{}, err
		}

		// Snapping will fail for locations that are over 300 meters away from a road, but at least
		// one should work since our sample size is large.
		snappedLocations, err := sv.SnapToRoads(randomLocations)
		if err != nil {
			return [][2]float64{}, err
		}

		// Except for when it fails anyway in subdivisions with a sparse road network (e.g., Roraima).
		if len(snappedLocations) == 1 && snappedLocations[0] == [2]float64{0, 0} {
			fmt.Println("no roads nearby")
			continue
		}

		// TODO: Why are there non-unique locations?
		uniqueLocations := make([][2]float64, 0)
		set := make(map[[2]float64]bool)
		for _, location := range snappedLocations {
			if set[location] {
				continue
			}
			uniqueLocations = append(uniqueLocations, location)
			set[location] = true
		}

		// There is no guarantee that valid Google Street View coverage exists at the snapped location.
		for _, location := range uniqueLocations {
			valid, err := sv.ValidateCoverage(location)
			if err != nil {
				return [][2]float64{}, err
			}
			if valid {
				locations = append(locations, location)
				fmt.Printf("found valid location %d\n", len(locations))
			} else {
				fmt.Println("nonexistent or invalid coverage at location")
			}
			if len(locations) == count {
				break
			}
		}
	}

	return locations, nil
}
