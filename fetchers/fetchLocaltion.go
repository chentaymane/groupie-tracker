package zone

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Locations represents the structure of the location data from the API
type Locations struct {
	Locations []string `json:"locations"`
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
