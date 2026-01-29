package zone

import (
	"encoding/json"
	"net/http"
)

// Locations represents the structure of the location data from the API
type Locations struct {
	Locations []string `json:"locations"`
}
type AllLocations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

var locationsCache map[int][]string

// FetchLocation retrieves the locations for a given artist ID from the external API
func FetchLocation(id int) ([]string, error) {
	locations, err := FetchAllLocations()
	if err != nil {
		return nil, err
	}

	return locations[id], nil
}

func FetchAllLocations() (map[int][]string, error) {
	if locationsCache != nil {
		return locationsCache, nil
	}

	url := "https://groupietrackers.herokuapp.com/api/locations"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data AllLocations
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// map artistID → locations
	result := make(map[int][]string)
	for _, item := range data.Index {
		result[item.ID] = item.Locations
	}

	locationsCache = result
	return result, nil
}

func Getallolocations() []string {
	locations, _ := FetchAllLocations()
	var allloct []string
	for _, location := range locations {
		for _, l := range location {
			if checkrepeat(allloct, l+" - location") {
				allloct = append(allloct, l+" - location")
			}
		}
	}
	return allloct
}

func checkrepeat(Any []string, l string) bool {
	for _, s := range Any {
		if s == l {
			return false
		}
	}
	return true
}
