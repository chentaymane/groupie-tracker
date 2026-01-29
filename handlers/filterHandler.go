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

	search := strings.TrimSpace(r.FormValue("searchLocation"))
	filtered := artists
	validFilter := false

	search = strings.ToLower(search) // case insensitive

	if search != "" {
		lastDash := strings.LastIndex(search, "-")

		if lastDash != -1 {
			query := strings.TrimSpace(search[:lastDash])
			filterType := strings.TrimSpace(search[lastDash+1:])

			switch filterType {
			case "location":
				filtered, err = FilterByLocation(artists, query)
				if err != nil {
					HandleError(w, http.StatusInternalServerError, "Internal Server Error")
					return
				}
				validFilter = true

			case "artist/band":
				filtered = zone.FilterByName(artists, query)
				validFilter = true

			case "member":
				filtered = zone.FilterByMember(artists, query)
				validFilter = true

			case "first album":
				filtered = zone.FilterByFirstAlbum(artists, query)
				validFilter = true

			case "creation date":
				filtered = zone.FilterByCreationDate(artists, query)
				validFilter = true
			}
		}
	}

	data := FilterViewData{
		Artists:         artists,
		LocationSearch:  search,
		Locations:       zone.Getallolocations(),
		ArtistNames:     zone.GetAllNameOfAtrtist(),
		MemberNames:     zone.GetAllMemberNames(),
		FirstAlbumDates: zone.GetAllFirstAlbumDates(),
		CreationDates:   zone.GetAllCreationDates(),
	}

	if validFilter {
		data.Artists = filtered
	} else {
		data.Artists = []zone.Artist{}
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
