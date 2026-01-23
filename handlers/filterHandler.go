package zone

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	zone "zone/fetchers"
)

type FilterViewData struct {
	Artists []zone.Artist

	CreationFrom int
	CreationTo   int
	FirstFrom    int
	FirstTo      int

	MembersChecked map[int]bool
	LocationSearch string
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

	creationFrom := parseInt(r.FormValue("creationFrom"), 1958)
	creationTo := parseInt(r.FormValue("creationTo"), 2026)
	firstFrom := parseInt(r.FormValue("firstFrom"), 1958)
	firstTo := parseInt(r.FormValue("firstTo"), 2026)
	membersChecked := r.Form["members"]
	locationSearch := r.FormValue("searchLocation")

	normalizeRange(&creationFrom, &creationTo)
	normalizeRange(&firstFrom, &firstTo)

	choseFilter := false
	filtered := []zone.Artist{}

	for _, a := range artists {
		if creationFrom != 1958 || creationTo != 2026 {
			choseFilter = true
			if a.CreationDate < creationFrom || a.CreationDate > creationTo {
				continue
			}
		}

		if firstFrom != 1958 || firstTo != 2026 {
			choseFilter = true
			year, _ := strconv.Atoi(a.FirstAlbum[len(a.FirstAlbum)-4:])
			if year < firstFrom || year > firstTo {
				continue
			}
		}

		if len(membersChecked) > 0 {
			choseFilter = true
			if !matchMembers(a, membersChecked) {
				continue
			}
		}

		filtered = append(filtered, a)
	}

	if locationSearch != "" {
		choseFilter = true
		filtered, err = FilterByLocation(filtered, locationSearch)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	membersMap := make(map[int]bool)
	for _, m := range membersChecked {
		n, _ := strconv.Atoi(m)
		membersMap[n] = true
	}

	data := FilterViewData{
		Artists:        artists, // default
		CreationFrom:   creationFrom,
		CreationTo:     creationTo,
		FirstFrom:      firstFrom,
		FirstTo:        firstTo,
		MembersChecked: membersMap,
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
