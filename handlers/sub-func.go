package zone

import (
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

func FilterByLocation(artists []zone.Artist, search string) ([]zone.Artist,) {
	allLocations, _ := zone.FetchAllLocations()


	search = strings.ToLower(search)
	var result []zone.Artist

	for _, artist := range artists {
		locations := allLocations[artist.ID]

		for _, loc := range locations {
			if strings.Contains(strings.ToLower(loc), search) {

				result = append(result, artist)
				break
			}
		}
	}

	return result
}
