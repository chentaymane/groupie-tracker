package zone

import "strings"

// FormatDate removes leading asterisks from date strings
func FormatDate(dates []string) []string {
	for i, d := range dates {
		dates[i] = strings.TrimPrefix(d, "*")
	}
	return dates
}
