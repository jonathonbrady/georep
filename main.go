package main

import (
	"flag"
	"fmt"
	"georep/geoguessr"
	"georep/overpass"
	"georep/streetview"
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var (
		country string
		road    string
		user    string
	)

	flag.StringVar(&country, "country", "", "country containing the road")
	flag.StringVar(&road, "road", "", "road within the country")
	flag.StringVar(&user, "user", "", "user id")

	flag.Parse()
	if country == "" || road == "" || user == "" {
		log.Fatalf("country, road, and user must be specified")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env file: %v", err)
	}

	gc, err := geoguessr.NewGeoguessrClient()
	if err != nil {
		log.Fatalf("creating geoguessr client: %v", err)
	}

	op, err := overpass.NewOverpassClient()
	if err != nil {
		log.Fatalf("creating overpass client: %v", err)
	}

	sv, err := streetview.NewStreetViewClient()
	if err != nil {
		log.Fatalf("creating google maps client: %v", err)
	}

	err = deleteOldMaps(gc, time.Duration(24 * time.Hour))
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

	latlongs, err := op.GetLocationsOnRoad(country, road)
	if err != nil {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("retrieving locations on %s in %s: %v\n", road, country, err)
	}
	log.Printf("retrieved locations on %s in %s\n", road, country)

	// TODO: Double check the distribution on this. It might skew towards the start/end of the road.
	rand.Shuffle(len(latlongs), func(i, j int) {
		latlongs[i], latlongs[j] = latlongs[j], latlongs[i]
	})

	locations := make([]geoguessr.Location, 0)
	for _, latlong := range latlongs {
		pass, err := sv.ValidateCoverage(latlong)
		if err != nil {
			delete := geoguessr.DeleteMapRequest{
				Id: mapId,
			}
			gc.DeleteMap(delete)
			log.Fatalf("validating coverage at %v: %v", latlong, err)
		}

		if !pass {
			log.Printf("invalid coverage at %v\n", latlong)
			continue
		}
		log.Printf("valid coverage at %v\n", latlong)

		location := geoguessr.Location{
			Heading:   0,
			Latitude:  latlong.Latitude,
			Longitude: latlong.Longitude,
			Pitch:     0,
			Zoom:      0,
		}
		locations = append(locations, location)

		if len(locations) == 5 {
			break
		}
	}

	if len(locations) != 5 {
		delete := geoguessr.DeleteMapRequest{
			Id: mapId,
		}
		gc.DeleteMap(delete)
		log.Fatalf("failed to find 5 locations")
	}

	update := geoguessr.UpdateMapRequest{
		Avatar: geoguessr.Avatar{
			Background: "day",
			Decoration: "cactus",
			Ground:     "green",
			Landscape:  "mountains",
		},
		Locations:   locations,
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
}

func deleteOldMaps(gc *geoguessr.GeoguessrClient, d time.Duration) error {
	maps, err := gc.ListMaps()
	if err != nil {
		return fmt.Errorf("listing maps: %v", err)
	}

	for _, m := range maps {
		if time.Now().Add(-d).After(m.CreatedAt) {
			delete := geoguessr.DeleteMapRequest{
				Id: m.ID,
			}
			err = gc.DeleteMap(delete)
			if err != nil {
				return fmt.Errorf("deleting map %s: %v", m.ID, err)
			}
		}
	}

	log.Printf("deleted %d maps\n", len(maps))
	return nil
}
