package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type FetchData struct {
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
}

func FetchSongs(externalAPI, group, song string) (*FetchData, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", externalAPI, group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to connect to external API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status code %d", resp.StatusCode)
	}

	var details FetchData
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, errors.New("failed to decode external API response")
	}

	return &details, nil

}
