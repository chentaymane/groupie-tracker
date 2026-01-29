package zone

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Artist represents the structure of an artist
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// FetchARtists retrieves the list of artists from the public API
func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

func GetAllNameOfAtrtist() []string {
	Artists, _ := FetchArtists()
	artistsNames := []string{}
	for _, r := range Artists {
		if checkrepeat(artistsNames, r.Name) {
			artistsNames = append(artistsNames, r.Name+" - member")
		}
	}
	return artistsNames
}

func FilterByName(artists []Artist,name string) ([]Artist) {
	

	var filtred []Artist
	for _, r := range artists {
		if strings.Contains(strings.ToLower(r.Name), strings.ToLower(name)) {
			filtred = append(filtred, r)
			break
		}
	}
	return filtred
}
