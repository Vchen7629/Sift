package service

import (
	"fmt"
	"time"
)

func FormatRelativeDate(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return "unknown"
	}

	days := int(time.Since(t).Hours() / 24)
	switch {
	case days < 1:
		h := int(time.Since(t).Hours())
		return pluralFormat(h, "hour")
	case days < 7:
		return pluralFormat(days, "day")
	case days < 30:
		return pluralFormat(days/7, "week")
	case days < 365:
		return pluralFormat(days/30, "month")
	default:
		return pluralFormat(days/365, "year")
	}
}

func pluralFormat(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s ago", unit)
	}

	return fmt.Sprintf("%d %ss ago", n, unit)
}