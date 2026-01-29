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
	Names          []string
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

	Search := r.FormValue("searchLocation")
	filtered := []zone.Artist{}
	choseFilter := false
	if Search != "" {
		if strings.HasSuffix(Search, " - location") {
			choseFilter = true
			Search = strings.TrimSuffix(Search, " - location")

			filtered, err = FilterByLocation(artists, Search)
			if err != nil {
				HandleError(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}

		} else if strings.HasSuffix(Search, " - member") {
			choseFilter = true
			Search = strings.TrimSuffix(Search, " - member")

			filtered = zone.FilterByName(artists, Search)

		}
	}

	data := FilterViewData{
		Artists: artists, // default

		LocationSearch: Search,
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
