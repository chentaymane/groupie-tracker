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

func FilterByLocation(artists []zone.Artist, alllocations *zone.AllLocations, search string) []zone.Artist {
	search = strings.ToLower(search)
	var result []zone.Artist

	for _, artist := range artists {
		for _, allloc := range alllocations.Index {
			if allloc.ID == artist.ID {
				for _, loc := range allloc.Locations {
					if strings.Contains(strings.ToLower(loc), search) {

						result = append(result, artist)
						break
					}
				}
				break
			}
		}
	}

	return result
}
