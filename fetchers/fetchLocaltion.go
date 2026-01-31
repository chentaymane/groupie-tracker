package zone

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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

var locationsCache *AllLocations

// FetchLocation retrieves the locations for a given artist ID from the external API
func FetchLocation(id int) ([]string, error) {
	locations, err := FetchAllLocations()
	if err != nil {
		return nil, err
	}

	for _, idx := range locations.Index {
		if idx.ID == id {
			return idx.Locations, nil
		}
	}

	return nil, errors.New("Not found")
}

func FetchAllLocations() (*AllLocations, error) {
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

	return &data, nil
}

func Getallolocations(locations *AllLocations) []string {
	var allloct []string
	for _, idx := range locations.Index {
		for _, l := range idx.Locations {
			if checkrepeat(allloct, l+" - location") {
				allloct = append(allloct, l+" - location")
			}
		}
	}
	return allloct
}

func FilterByLocation(artists []Artist, alllocations *AllLocations, search string) []Artist {
	search = strings.ToLower(search)
	var result []Artist

	for _, artist := range artists {
		for _, allloc := range alllocations.Index {
			if allloc.ID == artist.ID {
				for _, loc := range allloc.Locations {
					if strings.Contains(strings.ToLower(loc), search) {

						result = append(result, artist)
						break
					}
				}
				break
			}
		}
	}

	return result
}

func checkrepeat(Any []string, l string) bool {
	for _, s := range Any {
		if s == l {
			return false
		}
	}
	return true
}
