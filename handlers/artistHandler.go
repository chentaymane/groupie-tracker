package zone

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	zone "zone/fetchers"
)

// HandlerArtist serves the artist detail page
func HandlerArtist(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")

	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	if idStr == "" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, http.StatusNotFound, "Artist not found")
		return
	}

	artists, err := zone.FetchArtists()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	var artist zone.Artist

	found := false
	for _, a := range artists {
		if a.ID == id {
			artist = a
			found = true
			break
		}
	}

	if !found {
		HandleError(w, http.StatusNotFound, "Artist not found")
		return
	}

	locations, err := zone.FetchLocation(artist.ID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	dates, err := zone.FetchDate(artist.ID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	dates = zone.FormatDate(dates)
	relations, err := zone.FetchRelations(artist.ID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	data := struct {
		Artist    zone.Artist
		Relations map[string][]string
		Locations []string
		Dates     []string
	}{
		Artist:    artist,
		Relations: relations,
		Locations: locations,
		Dates:     dates,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
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
