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

func GetFilterSuggestions(artists []zone.Artist, locations *zone.AllLocations) []string {
	total := []string{}
	total = append(total, zone.Getallolocations(locations)...)
	total = append(total, zone.GetAllNameOfAtrtist(artists)...)
	total = append(total, zone.GetAllMemberNames(artists)...)
	total = append(total, zone.GetAllFirstAlbumDates(artists)...)
	total = append(total, zone.GetAllCreationDates(artists)...)
	return total
}
