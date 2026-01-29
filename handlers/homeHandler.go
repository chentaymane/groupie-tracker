package zone

import (
	"bytes"
	"html/template"
	"net/http"

	zone "zone/fetchers"
)

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	artists, err := zone.FetchArtists()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	data := FilterViewData{
		Artists: artists,

		Locations:       zone.Getallolocations(),
		ArtistNames:     zone.GetAllNameOfAtrtist(),
		MemberNames:     zone.GetAllMemberNames(),
		FirstAlbumDates: zone.GetAllFirstAlbumDates(),
		CreationDates:   zone.GetAllCreationDates(),
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	buf.WriteTo(w)
}
