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
	Locations      []string
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

	locationSearch := r.FormValue("searchLocation")
	if strings.HasSuffix(locationSearch, "- location") {
		locationSearch = strings.TrimSuffix(locationSearch, "- location")
	}
	filtered := []zone.Artist{}
	choseFilter := false

	if locationSearch != "" {
		choseFilter = true
		filtered, err = FilterByLocation(artists, locationSearch)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	data := FilterViewData{
		Artists: artists, // default

		LocationSearch: locationSearch,
	}

	if choseFilter {
		data.Artists = filtered
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
