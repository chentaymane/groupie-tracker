package zone

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"time"

	zone "zone/fetchers"
)

func HandleFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	creationFromStr := r.FormValue("creationFrom")
	creationToStr := r.FormValue("creationTo")

	firstFromStr := r.FormValue("firstFrom")
	firstToStr := r.FormValue("firstTo")

	creationFrom, err := strconv.Atoi(creationFromStr)
	if err != nil {
		creationFrom = 1958
	}
	creationTo, err := strconv.Atoi(creationToStr)
	if err != nil {
		creationTo = time.Now().Year()
	}
	firstFrom, err := strconv.Atoi(firstFromStr)
	if err != nil {
		firstFrom = 1958
	}
	firstTo, err := strconv.Atoi(firstToStr)
	if err != nil {
		firstTo = time.Now().Year()
	}
	if firstTo < firstFrom {
		temp := firstFrom
		firstFrom = firstTo
		firstTo = temp
	}

	if creationTo < creationFrom {
		temp := creationFrom
		creationFrom = creationTo
		creationTo = temp
	}

	membersCheckedStr := r.Form["members"] // slice of strings

	artists, err := zone.FetchArtists()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var ArtistsFiltred []zone.Artist
	choseRange := false
	if creationFrom != 1958 || creationTo != 2026 {
		choseRange = true
		for _, r := range artists {
			if r.CreationDate >= creationFrom && r.CreationDate <= creationTo && !containsArtistByID(ArtistsFiltred, r.ID) {
				ArtistsFiltred = append(ArtistsFiltred, r)
			}
		}
	}
	if firstFrom != 1958 || firstTo != 2026 {
		choseRange = true
		for _, r := range artists {
			year, _ := strconv.Atoi(r.FirstAlbum[len(r.FirstAlbum)-4:])
			if year >= firstFrom && year <= firstTo && !containsArtistByID(ArtistsFiltred, r.ID) {
				ArtistsFiltred = append(ArtistsFiltred, r)
			}
		}
	}
	if len(membersCheckedStr) != 0 {
		choseRange = true

		for _, m := range membersCheckedStr {
			n, err := strconv.Atoi(m)
			if err != nil {
				HandleError(w, http.StatusNotFound, "Page not found")
				return
			}
			for _, a := range artists {
				year, _ := strconv.Atoi(a.FirstAlbum[len(a.FirstAlbum)-4:])

				if len(a.Members) == n && a.CreationDate >= creationFrom && a.CreationDate <= creationTo && !containsArtistByID(ArtistsFiltred, a.ID) && year >= firstFrom && year <= firstTo {
					ArtistsFiltred = append(ArtistsFiltred, a)
				}
			}

		}
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var buf bytes.Buffer
	if !choseRange {
		if err := tmpl.Execute(&buf, artists); err != nil {
			HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
		buf.WriteTo(w)

	} else {

		if err := tmpl.Execute(&buf, ArtistsFiltred); err != nil {
			HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}

		buf.WriteTo(w)
	}
}

func containsArtistByID(list []zone.Artist, id int) bool {
	for _, a := range list {
		if a.ID == id {
			return true
		}
	}
	return false
}
