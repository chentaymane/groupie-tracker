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

	locations, err := zone.FetchAllLocations()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	search := strings.ToLower(strings.TrimSpace(r.FormValue("searchLocation")))

	var filtered []zone.Artist
	if search != "" {
		lastDash := strings.LastIndex(search, " - ")
		query := search
		queryType := ""
		if lastDash != -1 {
			query = strings.TrimSpace(search[:lastDash])
			queryType = strings.TrimSpace(search[lastDash+3:])
		}
		filtered = SearchEverywhere(artists, locations, query, queryType)
	} else {
		// no query, show all
		filtered = artists
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

func SearchEverywhere(artists []zone.Artist, locations *zone.AllLocations, query string, queryType string) []zone.Artist {
	if queryType == "" {
		// use all categories
		out := make(chan []zone.Artist, 5)

		go func() { out <- zone.FilterByName(artists, query) }()
		go func() { out <- zone.FilterByMember(artists, query) }()
		go func() { out <- FilterByLocation(artists, locations, query) }()
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
	} else {

		var filtered []zone.Artist
		switch queryType {
		case "location":
			filtered = FilterByLocation(artists, locations, query)

		case "artist/band":
			filtered = zone.FilterByName(artists, query)

		case "member":
			filtered = zone.FilterByMember(artists, query)

		case "first album":
			filtered = zone.FilterByFirstAlbum(artists, query)

		case "creation date":
			filtered = zone.FilterByCreationDate(artists, query)
		}

		return filtered
	}
}
