package datetime

import "time"

func FormatTime(t string) string {
	layouts := []string{"15:04:05", "15:04"}
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, t)
		if err == nil {
			return parsedTime.Format("3:04 PM")
		}
	}

	return ""
}

func FormatDate(t string) string {
	layouts := []string{time.RFC3339, "2006-01-02"}
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, t)
		if err == nil {
			return parsedTime.Format("January 2, 2006")
		}
	}

	return ""
}
