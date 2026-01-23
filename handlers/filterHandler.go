package zone

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"strings"

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
	tmpl.Execute(&buf, data)
	buf.WriteTo(w)
}

func parseInt(value string, def int) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return n
}

func normalizeRange(min, max *int) {
	if *max < *min {
		*min, *max = *max, *min
	}
}

func matchMembers(a zone.Artist, members []string) bool {
	for _, m := range members {
		n, _ := strconv.Atoi(m)
		if len(a.Members) == n {
			return true
		}
	}
	return false
}

func FilterByLocation(artists []zone.Artist, search string) ([]zone.Artist, error) {
	if search == "" {
		return artists, nil
	}

	allLocations, err := zone.FetchAllLocations()
	if err != nil {
		return nil, err
	}

	locMap := make(map[int][]string)
	for id, locs := range allLocations {
		locMap[id] = locs
	}

	search = strings.ToLower(search)
	var result []zone.Artist

	for _, artist := range artists {
		locations := locMap[artist.ID]
		for _, loc := range locations {
			if strings.Contains(strings.ToLower(loc), search) {
				result = append(result, artist)
				break
			}
		}
	}

	return result, nil
}
