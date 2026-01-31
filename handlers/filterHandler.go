package zone

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strings"

	zone "zone/fetchers"
)

type FilterViewData struct {
	Artists     []zone.Artist
	Search      string
	SearchError string
	Suggestions []string
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

	search := strings.ToLower(strings.TrimSpace(r.FormValue("search")))
	searchError := ""

	var filtered []zone.Artist
	if search != "" {
		lastDash := strings.LastIndex(search, " - ")
		query := search
		queryType := ""
		if lastDash != -1 {
			query = strings.TrimSpace(search[:lastDash])
			queryType = strings.TrimSpace(search[lastDash+3:])
		}
		filtered, err = SearchEverywhere(artists, locations, query, queryType)
		if err != nil {
			searchError = err.Error()
		}

	} else {
		// no query, show all
		filtered = artists
	}

	data := FilterViewData{
		Artists:     filtered,
		Search:      r.FormValue("search"),
		SearchError: searchError,
		Suggestions: GetFilterSuggestions(artists, locations),
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

func SearchEverywhere(artists []zone.Artist, locations *zone.AllLocations, query string, queryType string) ([]zone.Artist, error) {
	result := []zone.Artist{}
	if queryType == "" {
		// use all categories
		out := make(chan []zone.Artist, 5)

		go func() { out <- zone.FilterByName(artists, query) }()
		go func() { out <- zone.FilterByMember(artists, query) }()
		go func() { out <- zone.FilterByLocation(artists, locations, query) }()
		go func() { out <- zone.FilterByFirstAlbum(artists, query) }()
		go func() { out <- zone.FilterByCreationDate(artists, query) }()

		unique := make(map[int]zone.Artist)

		for i := 0; i < 5; i++ {
			results := <-out
			for _, a := range results {
				unique[a.ID] = a
			}
		}

		for _, a := range unique {
			result = append(result, a)
		}

	} else {
		switch queryType {
		case "location":
			result = zone.FilterByLocation(artists, locations, query)

		case "artist/band":
			result = zone.FilterByName(artists, query)

		case "member":
			result = zone.FilterByMember(artists, query)

		case "first album":
			result = zone.FilterByFirstAlbum(artists, query)

		case "creation date":
			result = zone.FilterByCreationDate(artists, query)
		default:
			return nil, errors.New("Bad category '" + queryType + "'")
		}
	}

	if len(result) == 0 {
		return nil, errors.New("Bad query '" + query + "'")
	}
	return result, nil
}
