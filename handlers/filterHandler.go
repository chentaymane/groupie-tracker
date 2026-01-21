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

	creationFrom, err := strconv.Atoi(creationFromStr)
	if err != nil {
		creationFrom = 1958
	}
	creationTo, err := strconv.Atoi(creationToStr)
	if err != nil {
		creationTo = time.Now().Year()
	}

	if creationTo < creationFrom {
		temp := creationFrom
		creationFrom = creationTo
		creationTo = temp
	}

	membersCheckedStr := r.Form["members"] // slice of strings

	var membersChecked []int
	for _, m := range membersCheckedStr {
		n, err := strconv.Atoi(m)
		if err == nil {
			membersChecked = append(membersChecked, n)
		}
	}
	artists, err := zone.FetchArtists()

	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	var ArtistsFiltred []zone.Artist
	for _, r := range artists {
		if len(r.Members) == len(membersChecked) && r.CreationDate >= creationFrom && r.CreationDate <= creationTo {
			ArtistsFiltred = append(ArtistsFiltred, r)
			continue
		}
		if r.CreationDate >= creationFrom && r.CreationDate <= creationTo {
			ArtistsFiltred = append(ArtistsFiltred, r)
		}

	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ArtistsFiltred); err != nil {
		HandleError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	buf.WriteTo(w)

}
