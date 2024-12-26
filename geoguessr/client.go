package geoguessr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type GeoguessrClient struct {
	client *http.Client
}

// Create a new challenge with these settings for a given map.
type CreateChallengeRequest struct {
	AccessLevel int    `json:"accessLevel"`
	NoMoving    bool   `json:"forbidMoving"`
	NoPanning   bool   `json:"forbidPanning"`
	NoZooming   bool   `json:"forbidZooming"`
	Map         string `json:"map"`
	TimeLimit   int    `json:"timeLimit"`
}

// Token is the challenge ID.
type CreateChallengeResponse struct {
	Token string `json:"token"`
}

// Create a new, unpopulated map.
type CreateMapRequest struct {
	Mode string `json:"mode"`
	Name string `json:"name"`
}

// Id is the ID of the new map.
type CreateMapResponse struct {
	Id string `json:"id"`
}

// Fabricated
type GetChallengeResultsRequest struct {
	Id string
}

type GetChallengeResultsResponse struct {
	Items []struct {
		GameToken  string `json:"gameToken"`
		PlayerName string `json:"playerName"`
		UserID     string `json:"userId"`
		TotalScore int    `json:"totalScore"`
		IsLeader   bool   `json:"isLeader"`
		PinURL     string `json:"pinUrl"`
		Game       struct {
			Token            string `json:"token"`
			Type             string `json:"type"`
			Mode             string `json:"mode"`
			State            string `json:"state"`
			RoundCount       int    `json:"roundCount"`
			TimeLimit        int    `json:"timeLimit"`
			ForbidMoving     bool   `json:"forbidMoving"`
			ForbidZooming    bool   `json:"forbidZooming"`
			ForbidRotating   bool   `json:"forbidRotating"`
			StreakType       string `json:"streakType"`
			Map              string `json:"map"`
			MapName          string `json:"mapName"`
			PanoramaProvider int    `json:"panoramaProvider"`
			Bounds           struct {
				Min struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"min"`
				Max struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"max"`
			} `json:"bounds"`
			Round  int `json:"round"`
			Rounds []struct {
				Lat                float64   `json:"lat"`
				Lng                float64   `json:"lng"`
				PanoID             string    `json:"panoId"`
				Heading            float64   `json:"heading"`
				Pitch              float64   `json:"pitch"`
				Zoom               float64   `json:"zoom"`
				StreakLocationCode string    `json:"streakLocationCode"`
				StartTime          time.Time `json:"startTime"`
			} `json:"rounds"`
			Player struct {
				TotalScore struct {
					Amount     string  `json:"amount"`
					Unit       string  `json:"unit"`
					Percentage float64 `json:"percentage"`
				} `json:"totalScore"`
				TotalDistance struct {
					Meters struct {
						Amount string `json:"amount"`
						Unit   string `json:"unit"`
					} `json:"meters"`
					Miles struct {
						Amount string `json:"amount"`
						Unit   string `json:"unit"`
					} `json:"miles"`
				} `json:"totalDistance"`
				TotalDistanceInMeters float64 `json:"totalDistanceInMeters"`
				TotalStepsCount       int     `json:"totalStepsCount"`
				TotalTime             int     `json:"totalTime"`
				TotalStreak           int     `json:"totalStreak"`
				Guesses               []struct {
					Lat               float64 `json:"lat"`
					Lng               float64 `json:"lng"`
					TimedOut          bool    `json:"timedOut"`
					TimedOutWithGuess bool    `json:"timedOutWithGuess"`
					SkippedRound      bool    `json:"skippedRound"`
					RoundScore        struct {
						Amount     string `json:"amount"`
						Unit       string `json:"unit"`
						Percentage int    `json:"percentage"`
					} `json:"roundScore"`
					RoundScoreInPercentage int `json:"roundScoreInPercentage"`
					RoundScoreInPoints     int `json:"roundScoreInPoints"`
					Distance               struct {
						Meters struct {
							Amount string `json:"amount"`
							Unit   string `json:"unit"`
						} `json:"meters"`
						Miles struct {
							Amount string `json:"amount"`
							Unit   string `json:"unit"`
						} `json:"miles"`
					} `json:"distance"`
					DistanceInMeters   float64 `json:"distanceInMeters"`
					StepsCount         int     `json:"stepsCount"`
					StreakLocationCode any     `json:"streakLocationCode"`
					Time               int     `json:"time"`
				} `json:"guesses"`
				IsLeader        bool `json:"isLeader"`
				CurrentPosition int  `json:"currentPosition"`
				Pin             struct {
					URL       string `json:"url"`
					Anchor    string `json:"anchor"`
					IsDefault bool   `json:"isDefault"`
				} `json:"pin"`
				NewBadges   []any  `json:"newBadges"`
				Explorer    any    `json:"explorer"`
				ID          string `json:"id"`
				Nick        string `json:"nick"`
				IsVerified  bool   `json:"isVerified"`
				Flair       int    `json:"flair"`
				CountryCode string `json:"countryCode"`
			} `json:"player"`
			ProgressChange struct {
				XpProgressions []struct {
					Xp           int `json:"xp"`
					CurrentLevel struct {
						Level   int `json:"level"`
						XpStart int `json:"xpStart"`
					} `json:"currentLevel"`
					NextLevel struct {
						Level   int `json:"level"`
						XpStart int `json:"xpStart"`
					} `json:"nextLevel"`
					CurrentTitle struct {
						ID           int    `json:"id"`
						TierID       int    `json:"tierId"`
						MinimumLevel int    `json:"minimumLevel"`
						Name         string `json:"name"`
					} `json:"currentTitle"`
				} `json:"xpProgressions"`
				AwardedXp struct {
					TotalAwardedXp int `json:"totalAwardedXp"`
					XpAwards       []struct {
						Xp     int    `json:"xp"`
						Reason string `json:"reason"`
						Count  int    `json:"count"`
					} `json:"xpAwards"`
				} `json:"awardedXp"`
				Medal                   int `json:"medal"`
				CompetitiveProgress     any `json:"competitiveProgress"`
				RankedSystemProgress    any `json:"rankedSystemProgress"`
				RankedTeamDuelsProgress any `json:"rankedTeamDuelsProgress"`
			} `json:"progressChange"`
		} `json:"game"`
	} `json:"items"`
	PaginationToken string `json:"paginationToken"`
}

// Fabricated
type PublishMapRequest struct {
	Id string
}

type PublishMapResponse struct {
	Message string `json:"message"`
}

type Avatar struct {
	Background string `json:"background"`
	Decoration string `json:"decoration"`
	Ground     string `json:"green"`
	Landscape  string `json:"landscape"`
}

type Location struct {
	Heading   float64 `json:"heading"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Pitch     float64 `json:"pitch"`
	Zoom      float64 `json:"zoom"`
}

// TODO: Maybe we'll want to support spaced repetition for regions, but that's probably hard and who cares.
// For now, all of the maps use coordinates to generate their locations.
type Region struct{}

// Add locations to a map.
type UpdateMapRequest struct {
	Avatar      Avatar     `json:"avatar"`
	Locations   []Location `json:"customCoordinates"`
	Description string     `json:"description"`
	Name        string     `json:"name"`
	Regions     []Region   `json:"regions"`
}

type UpdateMapResponse struct {
	Message string `json:"message"`
}

func NewGeoguessrClient(ncfa string) (*GeoguessrClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %v", err)
	}

	url := &url.URL{
		Scheme: "https",
		Host:   "www.geoguessr.com",
	}
	cookies := []*http.Cookie{
		{
			Name:   "_ncfa",
			Value:  ncfa,
			Domain: "www.geoguessr.com",
		},
	}
	jar.SetCookies(url, cookies)

	return &GeoguessrClient{
		client: &http.Client{Jar: jar},
	}, nil
}

// Generates a new challenge link for the requested map.
func (gc *GeoguessrClient) CreateChallenge(request CreateChallengeRequest) (string, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshaling request payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://www.geoguessr.com/api/v3/challenges", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := gc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status from challenges API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %v", err)
	}

	var response CreateChallengeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("unmarshaling response body: %v", err)
	}

	return fmt.Sprintf("https://geoguessr.com/maps/%s/play?challengeId=%s", request.Map, response.Token), nil
}

// Returns the map ID of the new map.
func (gc *GeoguessrClient) CreateMap(request CreateMapRequest) (string, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshaling request payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://www.geoguessr.com/api/v4/user-maps/drafts", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := gc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status from drafts API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %v", err)
	}

	var response CreateMapResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("unmarshaling response body: %v", err)
	}

	return response.Id, nil
}

func (gc *GeoguessrClient) GetChallengeResults(request GetChallengeResultsRequest) (*GetChallengeResultsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.geoguessr.com/api/v3/results/highscores/%s", request.Id), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status from drafts API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	var response *GetChallengeResultsResponse
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response body: %v", err)
	}

	return response, nil
}

func (gc *GeoguessrClient) PublishMap(request PublishMapRequest) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.geoguessr.com/api/v4/user-maps/drafts/%s/publish", request.Id), http.NoBody)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from drafts API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	var response PublishMapResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("unmarshaling response body: %v", err)
	}

	if response.Message != "OK" {
		return fmt.Errorf("bad message in response body: %v", response.Message)
	}

	return nil
}

func (gc *GeoguessrClient) UpdateMap(request UpdateMapRequest, id string) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshaling request payload: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.geoguessr.com/api/v4/user-maps/drafts/%s", id), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := gc.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	fmt.Println(string(payload))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from drafts API: %v ... %v", resp.StatusCode, string(body))
	}

	var response UpdateMapResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("unmarshaling response body: %v", err)
	}

	if response.Message != "OK" {
		return fmt.Errorf("bad message in response body: %v", response.Message)
	}

	return nil
}
