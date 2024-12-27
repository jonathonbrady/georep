package geoguessr

import (
	"net/http"
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
