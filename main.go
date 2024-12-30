package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"georep/data"
	"georep/geoguessr"
	"georep/googlemaps"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var (
		country     string
		road        string
		subdivision string
		user        string
	)

	flag.StringVar(&country, "country", "", "country containing the road")
	flag.StringVar(&road, "road", "", "road within the country")
	flag.StringVar(&subdivision, "subdivision", "", "first-order subdivision within the country")
	flag.StringVar(&user, "user", "", "user id")

	flag.Parse()
	if country == "" || user == "" {
		log.Fatalf("country and user must be specified")
	}
	if (road == "") == (subdivision == "") {
		log.Fatalf("either a road or first-order subdivision must be specified, but not both")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env file: %v", err)
	}

	gc, err := geoguessr.NewGeoguessrClient()
	if err != nil {
		log.Fatalf("creating geoguessr client: %v", err)
	}

	// TODO: Obsolete?
	// op, err := overpass.NewOverpassClient()
	// if err != nil {
	// 	log.Fatalf("creating overpass client: %v", err)
	// }

	sv, err := googlemaps.NewGoogleMapsClient()
	if err != nil {
		log.Fatalf("creating google maps client: %v", err)
	}

	err = deleteOldMaps(gc /* time.Duration(24 * time.Hour) */)
	if err != nil {
		log.Fatalf("deleting old maps: %v", err)
	}

	date := strings.Split(time.Now().Format(time.RFC3339), "T")[0]
	mapName := fmt.Sprintf("%s - %s", user, date)

	create := geoguessr.CreateMapRequest{
		Mode: "coordinates",
		Name: mapName,
	}
	mapId, err := gc.CreateMap(create)
	if err != nil {
		log.Fatalf("creating map: %v", err)
	}
	log.Printf(`created new map "%s" with id %s`, create.Name, mapId)

	locations := make([][2]float64, 0)
	if road != "" {
		// TODO: Use Google Directions API to get a polyline.
		// locations = generateLocationsOnRoad()
	} else {
		locations, err = data.GetLocationsInSubdivision(country, subdivision, 5, sv)
		if err != nil {
			log.Fatalf("getting locations in %v, %v: %v", subdivision, country, err)
		}
	}

	if len(locations) != 5 {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("failed to find 5 locations")
	}

	geoLocations := make([]geoguessr.Location, 0)
	for _, location := range locations {
		geoLocation := geoguessr.Location{
			Heading:   0,
			Latitude:  location[0],
			Longitude: location[1],
			Pitch:     0,
			Zoom:      0,
		}
		geoLocations = append(geoLocations, geoLocation)
	}

	update := geoguessr.UpdateMapRequest{
		Avatar: geoguessr.Avatar{
			Background: "day",
			Decoration: "cactus",
			Ground:     "green",
			Landscape:  "mountains",
		},
		Locations:   geoLocations,
		Description: "",
		Name:        mapName,
		Regions:     []geoguessr.Region{},
	}

	err = gc.UpdateMap(update, mapId)
	if err != nil {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("updating map %s: %v", mapId, err)
	}
	log.Printf("updated map %s with locations\n", mapId)

	publish := geoguessr.PublishMapRequest{
		Id: mapId,
	}
	err = gc.PublishMap(publish)
	if err != nil {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("publishing map %s: %v", mapId, err)
	}

	challenge := geoguessr.CreateChallengeRequest{
		AccessLevel: 1,
		NoMoving:    true,
		NoPanning:   false,
		NoZooming:   false,
		Map:         mapId,
		TimeLimit:   0,
	}
	link, err := gc.CreateChallenge(challenge)
	if err != nil {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("creating challenge for map %s: %v", mapId, err)
	}
	log.Println(link)

	calls := 0
	for _, n := range sv.APICalls {
		calls += n
	}
	fmt.Printf("used %d Google Maps API calls\n", calls)
	s, _ := json.MarshalIndent(sv.APICalls, "", "\t")
	fmt.Print(string(s))
}

func deleteOldMaps(gc *geoguessr.GeoguessrClient /* d time.Duration */) error {
	maps, err := gc.ListMaps()
	if err != nil {
		return fmt.Errorf("listing maps: %v", err)
	}

	deleted := 0
	for _, m := range maps {
		// if time.Now().Add(-d).After(m.CreatedAt) {
		delete := geoguessr.DeleteMapRequest{
			Id: m.ID,
		}
		err = gc.DeleteMap(delete)
		if err != nil {
			return fmt.Errorf("deleting map %s: %v", m.ID, err)
		}
		deleted++
		// }
	}

	log.Printf("deleted %d maps\n", deleted)
	return nil
}
