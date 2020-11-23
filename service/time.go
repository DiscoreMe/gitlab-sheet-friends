package service

import (
	"fmt"
	"time"
)

const timeLayout = "02.01.2006"

func TimeListName(t time.Time) string {
	f, l := firstWeekDay(t), lastWeekDay(t)
	return fmt.Sprintf("%s-%s", f.Format(timeLayout), l.Format(timeLayout))
}

func firstWeekDay(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		return t.Add(-6 * 24 * time.Hour)
	}
	return t.Add(time.Duration(1-int(t.Weekday())) * 24 * time.Hour)
}

func lastWeekDay(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		return t
	}
	return t.Add(time.Duration(7-int(t.Weekday())) * 24 * time.Hour)
}
