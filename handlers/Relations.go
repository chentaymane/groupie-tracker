package zone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Relations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func FetchRelations(id int) (map[string][]string, error) {
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
		return nil,err
	}

	return data.DatesLocations, nil
}
