package zone

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// FetchLocation retrieves the locations for a given artist ID from the external API
func FetchLocation(id int) ([]string, error) {
	relationsURL := "https://groupietrackers.herokuapp.com/api/locations/"
	resp, err := http.Get(relationsURL + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rel Locations
	err = json.NewDecoder(resp.Body).Decode(&rel)
	if err != nil {
		return nil, err
	}

	return rel.Locations, nil
}

func FetchAllLocations() (map[int][]string, error) {
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

	return result, nil
}

func Getallolocations(locations map[int][]string) []string {
	var allloct []string
	for _, location := range locations {
		for _, l := range location {
			if checkrepeat(allloct, l) {
				l = strings.ReplaceAll(l, "-", " ")
				allloct = append(allloct, l+"- location")
			}
		}
	}
	return allloct
}

func checkrepeat(locations []string, l string) bool {
	for _, s := range locations {
		if s == l {
			return false
		}
	}
	return true
}
