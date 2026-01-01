package zone

import "strings"


func FormatDate(dates []string) []string {
	for i, d := range dates {
		dates[i] = strings.TrimPrefix(d, "*")
	}
	return dates
}
