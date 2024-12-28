package geoguessr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func NewGeoguessrClient() (*GeoguessrClient, error) {
	ncfa, ok := os.LookupEnv("NCFA_COOKIE")
	if !ok {
		log.Fatal("ncfa cookie environment variable not set")
	}

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

func (gc *GeoguessrClient) DeleteMap(request DeleteMapRequest) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://www.geoguessr.com/api/v4/user-maps/%s", request.Id), http.NoBody)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from user-maps API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	var response DeleteMapResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("unmarshaling response body: %v", err)
	}

	if !response.Deleted {
		return fmt.Errorf("failed to delete map")
	}

	return nil
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

func (gc *GeoguessrClient) ListMaps() ([]Map, error) {
	req, err := http.NewRequest("GET", "https://www.geoguessr.com/api/v4/user-maps/maps", http.NoBody)
	if err != nil {
		return []Map{}, fmt.Errorf("creating request: %v", err)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return []Map{}, fmt.Errorf("executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Map{}, fmt.Errorf("bad status from user-maps API: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Map{}, fmt.Errorf("reading response body: %v", err)
	}

	var response []Map
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Map{}, fmt.Errorf("unmarshaling response body: %v", err)
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
