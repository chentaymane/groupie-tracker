package zone

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Dates represents the structure of the date data from the API
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

var datesCache []string

// FetchDate retrieves the dates for a given artist ID from the external API
func FetchDate(id int) ([]string, error) {
	if datesCache != nil {
		return datesCache, nil
	}

	url := "https://groupietrackers.herokuapp.com/api/dates/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data Dates
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	datesCache = data.Dates
	return data.Dates, nil
}
