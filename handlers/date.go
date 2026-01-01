package zone

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

func FetchDate(id int) ([]string, error) {
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

	return data.Dates, nil
}
