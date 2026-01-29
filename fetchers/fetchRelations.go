package zone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Relations represents the structure of the relations data from the API
type Relations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

var relationsCache map[string][]string

// FetchRelations retrieves the dates and locations for a given artist ID from the external API
func FetchRelations(id int) (map[string][]string, error) {
	if relationsCache != nil {
		return relationsCache, nil
	}

	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Err in get:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var data Relations
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	relationsCache = data.DatesLocations
	return data.DatesLocations, nil
}
