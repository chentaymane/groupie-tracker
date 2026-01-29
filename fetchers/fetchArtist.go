package zone

import (
	"encoding/json"
	"net/http"
	"strconv"
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
			artistsNames = append(artistsNames, r.Name+" - artist")
		}
	}
	return artistsNames
}

func FilterByName(artists []Artist, name string) []Artist {
	var filtred []Artist
	for _, r := range artists {
		if strings.Contains(strings.ToLower(r.Name), strings.ToLower(name)) {
			filtred = append(filtred, r)
			break
		}
	}
	return filtred
}

func GetAllMemberNames() []string {
	Artists, _ := FetchArtists()
	memberNames := []string{}
	for _, r := range Artists {
		for _, m := range r.Members {
			if checkrepeat(memberNames, m) {
				memberNames = append(memberNames, m+" - member")
			}
		}
	}
	return memberNames
}

func GetAllFirstAlbumDates() []string {
	Artists, _ := FetchArtists()
	dates := []string{}
	for _, r := range Artists {
		if checkrepeat(dates, r.FirstAlbum) {
			dates = append(dates, r.FirstAlbum+" - first album")
		}
	}
	return dates
}

func GetAllCreationDates() []string {
	Artists, _ := FetchArtists()
	dates := []string{}
	for _, r := range Artists {
		creationDateStr := strconv.Itoa(r.CreationDate)
		if checkrepeat(dates, creationDateStr) {
			dates = append(dates, creationDateStr+" - creation date")
		}
	}
	return dates
}

func FilterByMember(artists []Artist, member string) []Artist {
	member = strings.ToLower(member)
	var res []Artist

	for _, a := range artists {
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), member) {
				res = append(res, a)
				break
			}
		}
	}
	return res
}

func FilterByFirstAlbum(artists []Artist, year string) []Artist {
	var res []Artist
	for _, a := range artists {
		if strings.Contains(a.FirstAlbum, year) {
			res = append(res, a)
		}
	}
	return res
}

func FilterByCreationDate(artists []Artist, year string) []Artist {
	var res []Artist
	for _, a := range artists {
		if strconv.Itoa(a.CreationDate) == year {
			res = append(res, a)
		}
	}
	return res
}
