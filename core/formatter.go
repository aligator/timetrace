package core

import (
	"fmt"
	"time"
)

// Formatter represents a date- and time formatter. It provides all displayed
// date- and time layouts and is capable of parsing those layouts.
type Formatter struct {
	use12Hours bool
}

const dateLayout = "2006-01-02"

// ParseDate parses a date from an input string in the form YYYY-MM-DD. It also
// supports the `today` and `yesterday` aliases for convenience.
func (f *Formatter) ParseDate(input string) (time.Time, error) {
	if input == "today" {
		return time.Now(), nil
	}
	if input == "yesterday" {
		yesterday := time.Now().AddDate(0, 0, -1)
		return yesterday, nil
	}

	date, err := time.Parse(dateLayout, input)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// ParseTime parses a time from an input string in the configured timeLayout
func (f *Formatter) ParseTime(input string) (time.Time, error) {
	date, err := time.Parse(f.timeLayout(), input)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// CombineDateAndTime takes a date and a time and combines them to the time
// struct that represents the given time on the given day
func (f *Formatter) CombineDateAndTime(d, t time.Time) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
}

const (
	defaultTimeLayout        = "15:04"
	default12HoursTimeLayout = "03:04PM"
	defaultDatelayoutPretty  = "Monday, 02. January 2006"
)

func (f *Formatter) timeLayout() string {
	if f.use12Hours {
		return default12HoursTimeLayout
	}
	return defaultTimeLayout
}

func (f *Formatter) TimeString(input time.Time) string {
	return input.Format(f.timeLayout())
}

// PrettyDateString returns a nice representation of a given time
// example: Mon 31. May 2021
func (f *Formatter) PrettyDateString(input time.Time) string {
	return input.Format(defaultDatelayoutPretty)
}

const (
	defaultRecordKeyLayout        = "2006-01-02-15-04"
	default12HoursRecordKeyLayout = "2006-01-02-03-04PM"
)

func (f *Formatter) RecordKeyLayout() string {
	if f.use12Hours {
		return default12HoursRecordKeyLayout
	}
	return defaultRecordKeyLayout
}

// ParseRecordKey parses an input string in the form 2006-01-02-15-04 or
// 2006-01-02-03-04PM depending on the use12hours setting.
func (f *Formatter) ParseRecordKey(key string) (time.Time, error) {
	return time.Parse(f.RecordKeyLayout(), key)
}

func (f *Formatter) RecordKey(record *Record) string {
	return record.Start.Format(f.RecordKeyLayout())
}

// FormatTodayTime returns the formated string of the total
// time of today follwoing the format convention
func (f *Formatter) FormatTodayTime(report *Report) string {
	return f.FormatDuration(report.TrackedTimeToday)
}

// FormatCurrentTime returns the formated string of the current
// report time follwoing the format convention
func (f *Formatter) FormatCurrentTime(report *Report) string {
	return f.FormatDuration(*report.TrackedTimeCurrent)
}

// FormatBreakTime returns the formated string of the total time
// taking breaks today following the format convention
func (f *Formatter) FormatBreakTime(report *Report) string {
	return f.FormatDuration(report.BreakTimeToday)
}

// formatDuration formats the passed duration into a string.
// The format will be "8h 24min".
// seconds information is ignored.
func (f *Formatter) FormatDuration(duration time.Duration) string {

	hours := int64(duration.Hours()) % 60
	minutes := int64(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dmin", hours, minutes)
}
