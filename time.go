package timeapi

import (
	"fmt"
	"strconv"
	"time"
)

// Interval represents an interval (date and time) between two instants.
type Interval struct {
	year   int
	month  int
	day    int
	hour   int
	minute int
	second int
}

// NewInterval returns a new Interval instance.
func NewInterval(year, month, day, hour, minute, second int) Interval {
	return Interval{
		year:   year,
		month:  month,
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}
}

// NewIntervalDate returns a new Interval instance with only date component.
func NewIntervalDate(year, month, day int) Interval {
	return Interval{
		year:  year,
		month: month,
		day:   day,
	}
}

// NewIntervalTime returns a new Interval instance with only time component.
func NewIntervalTime(hour, minute, second int) Interval {
	return Interval{
		hour:   hour,
		minute: minute,
		second: second,
	}
}

func (i Interval) String() string {
	// special case if all values are zero
	if i.IsZero() {
		return "0"
	}

	var s string
	if i.year != 0 {
		s += strconv.Itoa(i.year) + "y"
	}
	if i.month != 0 {
		s += strconv.Itoa(i.month) + "mo"
	}
	if i.day != 0 {
		s += strconv.Itoa(i.day) + "d"
	}
	if i.hour != 0 {
		s += strconv.Itoa(i.hour) + "h"
	}
	if i.minute != 0 {
		s += strconv.Itoa(i.minute) + "m"
	}
	if i.second != 0 {
		s += strconv.Itoa(i.second) + "s"
	}
	return s
}

func (i Interval) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.String() + `"`), nil
}

func (i *Interval) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return NewErrJsonValue(fmt.Errorf("interval %q is invalid", string(b)))
	}

	b = b[1 : len(b)-1]
	in, err := ParseInterval(string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}
	*i = in
	return nil
}

// Is Zero reports whether i represents the zero interval.
func (i Interval) IsZero() bool {
	return i.year == 0 && i.month == 0 &&
		i.day == 0 && i.hour == 0 &&
		i.minute == 0 && i.second == 0
}

// Date returns the date component of the interval.
func (i Interval) Date() (year, month, day int) {
	return i.year, i.month, i.day
}

// Time returns the time component of the interval.
func (i Interval) Time() (hour, min, sec int) {
	return i.hour, i.minute, i.second
}

// Duration represents the duration between two instants.
type Duration struct {
	neg    int
	hour   int
	minute int
	second int
}

// NewDuration returns a new Duration instance.
func NewDuration(hour, minute, second int) Duration {
	neg := 1
	if hour < 0 || minute < 0 || second < 0 {
		neg = -1
	}
	return Duration{
		neg:    neg,
		hour:   abs(hour),
		minute: abs(minute),
		second: abs(second),
	}
}

func (d Duration) String() string {
	// special case if all values are zero
	// return 0h0m0s for better readability
	if d.hour == 0 && d.minute == 0 && d.second == 0 {
		return "0h0m0s"
	}

	var s string
	if d.neg == -1 {
		s += "-"
	}
	if d.hour != 0 {
		s += strconv.Itoa(d.hour) + "h"
	}
	if d.minute != 0 {
		s += strconv.Itoa(d.minute) + "m"
	}
	if d.second != 0 {
		s += strconv.Itoa(d.second) + "s"
	}
	return s
}

// IsZero reports whether d represents the zero duration.
func (d Duration) IsZero() bool {
	return d.hour == 0 && d.minute == 0 && d.second == 0
}

// Hours returns the duration hours.
func (d Duration) Hours() int {
	return d.neg * d.hour
}

// Minutes returns the duration minutes.
func (d Duration) Minutes() int {
	return d.neg * d.minute
}

// Seconds returns the duration seconds.
func (d Duration) Seconds() int {
	return d.neg * d.second
}

// GoDuration returns the standard go time.Duration instance.
func (d Duration) GoDuration() time.Duration {
	return time.Duration(d.neg)*time.Duration(d.hour)*time.Hour +
		time.Duration(d.minute)*time.Minute +
		time.Duration(d.second)*time.Second
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	// strip quotes
	if len(b) < 2 {
		return NewErrJsonValue(fmt.Errorf("duration %q is invalid", string(b)))
	}
	b = b[1 : len(b)-1]

	dur, err := ParseDuration(string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}
	*d = dur
	return nil
}

// Timezone represents a time zone.
type Timezone struct {
	loc time.Location
}

// NewWeekday returns a new Timezone instance.
func NewTimezone(loc time.Location) Timezone {
	return Timezone{loc: loc}
}

func (t Timezone) String() string {
	return t.loc.String()
}

// GoLocation returns the standard go time.Location instance.
func (t Timezone) GoLocation() *time.Location {
	return &t.loc
}

func (t Timezone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *Timezone) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return NewErrJsonValue(fmt.Errorf("timezone %q is invalid", string(b)))
	}
	b = b[1 : len(b)-1]

	if len(b) == 0 || string(b) == "Local" {
		return NewErrJsonValue(fmt.Errorf("timezone %q is invalid", string(b)))
	}

	loc, err := time.LoadLocation(string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}

	t.loc = *loc
	return nil
}

// Weekday represents a day of the week.
type Weekday struct {
	w time.Weekday
}

var weekdayNames = map[time.Weekday]string{
	time.Sunday:    "SUNDAY",
	time.Monday:    "MONDAY",
	time.Tuesday:   "TUESDAY",
	time.Wednesday: "WEDNESDAY",
	time.Thursday:  "THURSDAY",
	time.Friday:    "FRIDAY",
	time.Saturday:  "SATURDAY",
}

var namesToWeekday = map[string]time.Weekday{
	"SUNDAY":    time.Sunday,
	"MONDAY":    time.Monday,
	"TUESDAY":   time.Tuesday,
	"WEDNESDAY": time.Wednesday,
	"THURSDAY":  time.Thursday,
	"FRIDAY":    time.Friday,
	"SATURDAY":  time.Saturday,
}

// NewWeekday returns a new Weekday instance. It panics if the weekday is out of range.
func NewWeekday(w time.Weekday) Weekday {
	if w < time.Sunday || w > time.Saturday {
		panic(fmt.Sprintf("weekday %d is out of range", w))
	}
	return Weekday{w: w}
}

func (w Weekday) String() string {
	return weekdayNames[w.w]
}

// GoWeekday returns the standard go time.Weekday instance.
func (w Weekday) GoWeekday() time.Weekday {
	return w.w
}

func (w Weekday) MarshalJSON() ([]byte, error) {
	return []byte(`"` + w.String() + `"`), nil
}

func (w *Weekday) UnmarshalJSON(b []byte) error {
	// strip quotes
	if len(b) < 2 {
		return NewErrJsonValue(fmt.Errorf("weekday %q is invalid", string(b)))
	}
	b = b[1 : len(b)-1]

	weekday, ok := namesToWeekday[string(b)]
	if !ok {
		return NewErrJsonValue(fmt.Errorf("weekday invalid value %q", string(b)))
	}
	w.w = weekday
	return nil
}

// time layout
const (
	timeLayout       = "15:04:05"
	quotedTimeLayout = `"` + timeLayout + `"`
)

// Time represents a time (hour, minute, second) with UTC timezone.
type Time struct {
	hour int
	min  int
	sec  int
}

// NewTime returns a new Time instance.
// It panics if the hour, minute, or second is out of range.
func NewTime(hour, min, sec int) Time {
	if hour < 0 || hour > 23 {
		panic(fmt.Sprintf("hour %d is out of range", hour))
	}
	if min < 0 || min > 59 {
		panic(fmt.Sprintf("minute %d is out of range", min))
	}
	if sec < 0 || sec > 59 {
		panic(fmt.Sprintf("second %d is out of range", sec))
	}
	return Time{hour, min, sec}
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.hour, t.min, t.sec)
}

// Clock returns the hour, minute, and second
// within the day specified by t.
func (t Time) Clock() (hour, min, sec int) {
	return t.hour, t.min, t.sec
}

// IsZero reports whether t represents the zero time instant.
func (t Time) IsZero() bool {
	return t.hour == 0 && t.min == 0 && t.sec == 0
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool {
	return t.hour < u.hour ||
		(t.hour == u.hour && (t.min < u.min || (t.min == u.min && t.sec < u.sec)))
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	return t.hour > u.hour ||
		(t.hour == u.hour && (t.min > u.min || (t.min == u.min && t.sec > u.sec)))
}

// Equal reports whether t and u represent the same time instant.
func (t Time) Equal(u Time) bool {
	return t.hour == u.hour && t.min == u.min && t.sec == u.sec
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(quotedTimeLayout, string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}
	t.hour = tm.Hour()
	t.min = tm.Minute()
	t.sec = tm.Second()
	return nil
}

// date layout
const (
	dateLayout       = "2006-01-02"
	quotedDateLayout = `"` + dateLayout + `"`
)

// Date represents a date (year, month, day) with UTC timezone.
type Date struct {
	year  int
	month time.Month
	day   int
}

// NewDate returns a new Date instance.
// It panics if the month or day is out of their usual ranges.
func NewDate(year int, month time.Month, day int) Date {
	s := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	if _, err := time.Parse(dateLayout, s); err != nil {
		panic(fmt.Sprintf("date %04d-%02d-%02d: %s", year, month, day, err))
	}
	return Date{year, month, day}
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

// Date returns the year, month, and day in which d occurs.
func (d Date) Date() (year int, month time.Month, day int) {
	return d.year, d.month, d.day
}

// Before reports whether the date d is before u.
func (d Date) Before(u Date) bool {
	return d.year < u.year ||
		(d.year == u.year &&
			(d.month < u.month || (d.month == u.month && d.day < u.day)))
}

// After reports whether the date d is after u.
func (d Date) After(u Date) bool {
	return d.year > u.year ||
		(d.year == u.year &&
			(d.month > u.month || (d.month == u.month && d.day > u.day)))
}

// Equal reports whether d and u represent the same date.
func (d Date) Equal(u Date) bool {
	return d.year == u.year &&
		d.month == u.month &&
		d.day == u.day
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(quotedDateLayout, string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}

	d.year = tm.Year()
	d.month = tm.Month()
	d.day = tm.Day()
	return nil
}

// dateTime layout
const (
	dateTimeLayout       = "2006-01-02T15:04:05Z"
	quotedDateTimeLayout = `"` + dateTimeLayout + `"`
)

// DateTime represents a date and time with UTC timezone.
type DateTime struct {
	t time.Time
}

// NewDateTime returns a new DateTime instance. It panics if the month, day,
// hour, minute, or second is out of range.
// This is to prevent the zero date from being used.
func NewDateTime(year int, month time.Month, day, hour, min, sec int) DateTime {
	NewDate(year, month, day)
	NewTime(hour, min, sec)
	return DateTime{time.Date(year, month, day, hour, min, sec, 0, time.UTC)}
}

func (dt DateTime) String() string {
	return dt.t.Format(dateTimeLayout)
}

// Date returns the year, month, and day in which dt occurs.
func (dt DateTime) Date() (year int, month time.Month, day int) {
	return dt.t.Year(), dt.t.Month(), dt.t.Day()
}

// Clock returns the hour, minute, and second within the day specified by dt.
func (dt DateTime) Clock() (hour, min, sec int) {
	return dt.t.Hour(), dt.t.Minute(), dt.t.Second()
}

// GoTime returns the standard go time.Time instance.
func (dt DateTime) GoTime() time.Time {
	return dt.t
}

// Before reports whether the date and time dt is before u.
func (dt DateTime) Before(u DateTime) bool {
	return dt.t.Before(u.t)
}

// After reports whether the date and time dt is after u.
func (dt DateTime) After(u DateTime) bool {
	return dt.t.After(u.t)
}

// Equal reports whether dt and u represent the same date and time.
func (dt DateTime) Equal(u DateTime) bool {
	return dt.t.Equal(u.t)
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dt.String() + `"`), nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(quotedDateTimeLayout, string(b))
	if err != nil {
		return NewErrJsonValue(err)
	}
	dt.t = tm
	return nil
}
