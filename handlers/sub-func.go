package zone

import (
	"strconv"
	"strings"
	zone "zone/fetchers"
)

// FormatDate removes leading asterisks from date strings
func FormatDate(dates []string) []string {
	for i, d := range dates {
		dates[i] = strings.TrimPrefix(d, "*")
	}
	return dates
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

	search = strings.ToLower(search)
	var result []zone.Artist

	for _, artist := range artists {
		locations := allLocations[artist.ID]
		
		for _, loc := range locations {
			loc = strings.ReplaceAll(loc,"-"," ")
			if strings.Contains(strings.ToLower(loc), search) {
				result = append(result, artist)
				break
			}
		}
	}

	return result, nil
}
