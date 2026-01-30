package zone

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	zone "zone/fetchers"
)

type FilterViewData struct {
	Artists        []zone.Artist
	LocationSearch string

	Locations       []string
	ArtistNames     []string
	MemberNames     []string
	FirstAlbumDates []string
	CreationDates   []string
}

func HandleFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	artists, err := zone.FetchArtists()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	search := strings.ToLower(strings.TrimSpace(r.FormValue("searchLocation")))

	var filtered []zone.Artist
	if search != "" {
		lastDash := strings.LastIndex(search, " - ")
		if lastDash != -1 {
			search = strings.TrimSpace(search[:lastDash]) // remove type substring
		}

		filtered = SearchEverywhere(artists, search)
	}

	data := FilterViewData{
		Artists:         filtered,
		LocationSearch:  search,
		Locations:       zone.Getallolocations(),
		ArtistNames:     zone.GetAllNameOfAtrtist(),
		MemberNames:     zone.GetAllMemberNames(),
		FirstAlbumDates: zone.GetAllFirstAlbumDates(),
		CreationDates:   zone.GetAllCreationDates(),
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	buf.WriteTo(w)
}

func SearchEverywhere(artists []zone.Artist, query string) []zone.Artist {
	out := make(chan []zone.Artist, 5)

	go func() { out <- zone.FilterByName(artists, query) }()
	go func() { out <- zone.FilterByMember(artists, query) }()
	go func() { out <- FilterByLocation(artists, query) }()
	go func() { out <- zone.FilterByFirstAlbum(artists, query) }()
	go func() { out <- zone.FilterByCreationDate(artists, query) }()

	unique := make(map[int]zone.Artist)

	for i := 0; i < 5; i++ {
		results := <-out
		for _, a := range results {
			unique[a.ID] = a
		}
	}

	var merged []zone.Artist
	for _, a := range unique {
		merged = append(merged, a)
	}

	return merged
}
