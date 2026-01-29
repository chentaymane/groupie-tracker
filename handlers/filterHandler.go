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

	Search := r.FormValue("searchLocation")
	filtered := []zone.Artist{}
	validFilter := true

	// validate filter
	for {
		if Search != "" {
			validFilter = false
			break
		}

		split := strings.Split(Search, "-")
		if len(split) != 2 {
			validFilter = false
			break
		}
		query := strings.TrimSpace(split[0])
		queryType := strings.TrimSpace(split[1])

		switch queryType {
		case "location":
			filtered, err = FilterByLocation(artists, query)
			if err != nil {
				HandleError(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}

		case "artist":
			filtered = zone.FilterByName(artists, query)

		case "member":
		case "first album":
		case "creation date":
			// TODO...

		default:
			// invalid type
			validFilter = false
		}

		break
	}

	data := FilterViewData{
		Artists: artists, // default

		LocationSearch: Search,
	}

	if validFilter {
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
