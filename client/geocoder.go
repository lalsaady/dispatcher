package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/lalsaady/dispatcher/model"
)

type GeocoderClient interface {
	GetCoords(address string) (model.Points, error)
}

type Geocoder struct{}

func NewGeocoderClient() GeocoderClient {
	return &Geocoder{}
}

func (g *Geocoder) GetCoords(address string) (model.Points, error) {
	// URL encode the address
	encodedAddress := url.QueryEscape(address)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		encodedAddress,
		"key")

	resp, err := http.Get(url)
	if err != nil {
		return model.Points{}, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return model.Points{}, fmt.Errorf("failed to make request, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Points{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return model.Points{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result) > 0 {
		points, err := parseResult(result)
		if err != nil {
			return model.Points{}, fmt.Errorf("failed to parse result: %v", err)
		}
		return points, nil
	}

	return model.Points{}, fmt.Errorf("error getting coordinates from google maps api")
}

// Parse the result from API response
func parseResult(result map[string]any) (model.Points, error) {
	results, ok := result["results"].([]any)
	if !ok || len(results) == 0 {
		return model.Points{}, fmt.Errorf("no results found in response")
	}

	firstResult, ok := results[0].(map[string]any)
	if !ok {
		return model.Points{}, fmt.Errorf("invalid result format")
	}

	geometry, ok := firstResult["geometry"].(map[string]any)
	if !ok {
		return model.Points{}, fmt.Errorf("geometry not found")
	}

	location, ok := geometry["location"].(map[string]any)
	if !ok {
		return model.Points{}, fmt.Errorf("location not found")
	}

	lat, ok := location["lat"].(float64)
	if !ok {
		return model.Points{}, fmt.Errorf("latitude not found or invalid")
	}

	lng, ok := location["lng"].(float64)
	if !ok {
		return model.Points{}, fmt.Errorf("longitude not found or invalid")
	}

	return model.Points{
		Lat: lat,
		Lon: lng,
	}, nil
}
