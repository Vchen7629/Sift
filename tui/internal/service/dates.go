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

	const (
		oneDayHours   = 24
		oneWeekHours  = 168
		oneMonthHours = 720
		oneYearHours  = 8760
	)

	elapsed := time.Since(t)
	hours := int(elapsed.Hours())

	switch {
	case hours < 1:
		m := int(elapsed.Minutes())
		return pluralFormat(m, "minute")
	case hours < oneDayHours:
		return pluralFormat(hours, "hour")
	case hours < oneWeekHours:
		return pluralFormat(hours/oneDayHours, "day")
	case hours < oneMonthHours:
		return pluralFormat(hours/oneWeekHours, "week")
	case hours < oneYearHours:
		return pluralFormat(hours/oneMonthHours, "month")
	default:
		return pluralFormat(hours/oneYearHours, "year")
	}
}

func pluralFormat(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s ago", unit)
	}

	return fmt.Sprintf("%d %ss ago", n, unit)
}
