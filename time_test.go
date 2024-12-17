package timeapi

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/krhubert/assert"
)

func TestInterval(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewInterval(0, 0, 0, 0, 0, 0).String(), "0")
		assert.Equal(t, NewInterval(1, 2, 3, 4, 5, 6).String(), "1y2mo3d4h5m6s")
		assert.Equal(t, NewIntervalTime(4, 5, 6).String(), "4h5m6s")
		assert.Equal(t, NewIntervalDate(4, 5, 6).String(), "4y5mo6d")
	})

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, NewInterval(0, 0, 0, 0, 0, 0).IsZero())
		assert.False(t, NewInterval(0, 0, 0, 0, 0, 1).IsZero())
	})

	t.Run("Date", func(t *testing.T) {
		year, month, day := NewInterval(1, 2, 3, 4, 5, 6).Date()
		assert.Equal(t, year, 1)
		assert.Equal(t, month, 2)
		assert.Equal(t, day, 3)
	})

	t.Run("Time", func(t *testing.T) {
		hour, min, sec := NewInterval(1, 2, 3, 4, 5, 6).Time()
		assert.Equal(t, hour, 4)
		assert.Equal(t, min, 5)
		assert.Equal(t, sec, 6)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		i := NewInterval(1, 2, 3, 4, 5, 6)
		out, err := json.Marshal(i)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"1y2mo3d4h5m6s"`)

		i = NewIntervalTime(4, 5, 6)
		out, err = json.Marshal(i)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"4h5m6s"`)

		i = NewIntervalDate(4, 5, 6)
		out, err = json.Marshal(i)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"4y5mo6d"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var i Interval
		err := json.Unmarshal([]byte(`""`), &i)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`"0"`), &i)
		assert.NoError(t, err)

		err = json.Unmarshal([]byte(`"1y2mo3d4h5m6s"`), &i)
		assert.NoError(t, err)
		assert.Equal(t, i, NewInterval(1, 2, 3, 4, 5, 6))

		err = json.Unmarshal([]byte(`1`), &i)
		assert.ErrorContains(t, err, "is invalid")

		err = json.Unmarshal([]byte(`"."`), &i)
		assert.ErrorContains(t, err, "invalid interval")

		err = json.Unmarshal([]byte(`"1sm"`), &i)
		assert.ErrorContains(t, err, "unknown unit \"sm\"")

		err = json.Unmarshal([]byte(`"1s1s"`), &i)
		assert.ErrorContains(t, err, "unit \"s\" repeated")

		err = json.Unmarshal([]byte(`"1s1m"`), &i)
		assert.ErrorContains(t, err, "unit \"m\" must be in the order")

		err = json.Unmarshal([]byte(`"1y1mo4h3d"`), &i)
		assert.ErrorContains(t, err, "unit \"d\" must be in the order")

		err = json.Unmarshal([]byte(`"1"`), &i)
		assert.ErrorContains(t, err, "missing unit")
	})
}

func TestDuration(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewDuration(0, 0, 0).String(), "0h0m0s")
		assert.Equal(t, NewDuration(0, 0, 1).String(), "1s")
		assert.Equal(t, NewDuration(0, 1, 0).String(), "1m")
		assert.Equal(t, NewDuration(1, 0, 0).String(), "1h")
		assert.Equal(t, NewDuration(1, 2, 3).String(), "1h2m3s")
	})

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, NewDuration(0, 0, 0).IsZero())
		assert.False(t, NewDuration(0, 0, 1).IsZero())
		assert.False(t, NewDuration(1, 0, 0).IsZero())
		assert.False(t, NewDuration(0, 1, 0).IsZero())
	})

	t.Run("Hours", func(t *testing.T) {
		assert.Equal(t, NewDuration(0, 0, 0).Hours(), 0)
		assert.Equal(t, NewDuration(1, 0, 0).Hours(), 1)
		assert.Equal(t, NewDuration(2, 70, 59).Hours(), 2)
	})

	t.Run("Minutes", func(t *testing.T) {
		assert.Equal(t, NewDuration(0, 0, 0).Minutes(), 0)
		assert.Equal(t, NewDuration(0, 1, 61).Minutes(), 1)
		assert.Equal(t, NewDuration(3, 1, 61).Minutes(), 1)
	})

	t.Run("Seconds", func(t *testing.T) {
		assert.Equal(t, NewDuration(0, 0, 0).Seconds(), 0)
		assert.Equal(t, NewDuration(0, 0, 59).Seconds(), 59)
		assert.Equal(t, NewDuration(0, 0, 61).Seconds(), 61)
		assert.Equal(t, NewDuration(1, 2, 61).Seconds(), 61)
	})

	t.Run("GoDuration", func(t *testing.T) {
		assert.Equal(t, NewDuration(0, 0, 0).GoDuration(), 0)
		assert.Equal(t, NewDuration(1, 2, 3).GoDuration(), time.Hour+2*time.Minute+3*time.Second)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		d := NewDuration(0, 0, 0)
		out, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"0h0m0s"`)

		d = NewDuration(0, 2, 3)
		out, err = json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"2m3s"`)

		d = NewDuration(1, 2, 3)
		out, err = json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"1h2m3s"`)

		d = NewDuration(-1, 2, 3)
		out, err = json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"-1h2m3s"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var d Duration
		err := json.Unmarshal([]byte(`"0h0m0s"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDuration(0, 0, 0))

		err = json.Unmarshal([]byte(`"0"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDuration(0, 0, 0))

		err = json.Unmarshal([]byte(`"1h2m3s"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDuration(1, 2, 3))

		err = json.Unmarshal([]byte(`"-1h2m3s"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDuration(-1, 2, 3))

		err = json.Unmarshal([]byte(`"2h59m59s"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDuration(2, 59, 59))

		err = json.Unmarshal([]byte(`""`), &d)
		assert.ErrorContains(t, err, "invalid duration")

		err = json.Unmarshal([]byte(`"1"`), &d)
		assert.ErrorContains(t, err, "missing unit")

		err = json.Unmarshal([]byte(`"."`), &d)
		assert.ErrorContains(t, err, "invalid duration")

		err = json.Unmarshal([]byte(`"1hh10ns"`), &d)
		assert.ErrorContains(t, err, "unknown unit \"hh\"")

		err = json.Unmarshal([]byte(`"1h10ns"`), &d)
		assert.ErrorContains(t, err, "unknown unit \"ns\"")

		err = json.Unmarshal([]byte(`"1h10us"`), &d)
		assert.ErrorContains(t, err, "unknown unit \"us\"")

		err = json.Unmarshal([]byte(`"1h10ms"`), &d)
		assert.ErrorContains(t, err, "unknown unit \"ms\"")

		err = json.Unmarshal([]byte(`"1h1h"`), &d)
		assert.ErrorContains(t, err, "unit \"h\" repeated")

		err = json.Unmarshal([]byte(`"1s1h"`), &d)
		assert.ErrorContains(t, err, "unit \"h\" must be in the order")

		err = json.Unmarshal([]byte(`"1m1h"`), &d)
		assert.ErrorContains(t, err, "unit \"h\" must be in the order")

		err = json.Unmarshal([]byte(`"1s1m"`), &d)
		assert.ErrorContains(t, err, "unit \"m\" must be in the order")

		err = json.Unmarshal([]byte(`0`), &d)
		assert.ErrorContains(t, err, "is invalid")
	})
}

func TestTimezone(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewTimezone(*time.UTC).String(), "UTC")
		assert.Equal(t, NewTimezone(*loc).String(), "America/New_York")
	})

	t.Run("GoLocation", func(t *testing.T) {
		assert.Equal(t, NewTimezone(*time.UTC).GoLocation(), time.UTC)
		assert.Equal(t, NewTimezone(*loc).GoLocation(), loc)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		tz := NewTimezone(*time.UTC)
		out, err := json.Marshal(tz)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"UTC"`)

		tz = NewTimezone(*loc)
		out, err = json.Marshal(tz)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"America/New_York"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var tz Timezone
		err := json.Unmarshal([]byte(`"UTC"`), &tz)
		assert.NoError(t, err)
		assert.Equal(t, tz, NewTimezone(*time.UTC))

		err = json.Unmarshal([]byte(`"America/New_York"`), &tz)
		assert.NoError(t, err)
		assert.Equal(t, tz, NewTimezone(*loc))

		err = json.Unmarshal([]byte(`1`), &tz)
		assert.ErrorContains(t, err, "is invalid")

		err = json.Unmarshal([]byte(`"Local"`), &tz)
		assert.ErrorContains(t, err, "is invalid")

		err = json.Unmarshal([]byte(`"America/Invalid"`), &tz)
		assert.ErrorContains(t, err, "unknown")

		err = json.Unmarshal([]byte(`""`), &tz)
		assert.ErrorContains(t, err, "is invalid")
	})
}

func TestWeekday(t *testing.T) {
	t.Run("NewWeekday", func(t *testing.T) {
		assert.Panic(t, func() { NewWeekday(time.Sunday - 1) })
		assert.Panic(t, func() { NewWeekday(time.Saturday + 1) })
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewWeekday(time.Sunday).String(), "SUNDAY")
		assert.Equal(t, NewWeekday(time.Monday).String(), "MONDAY")
		assert.Equal(t, NewWeekday(time.Tuesday).String(), "TUESDAY")
		assert.Equal(t, NewWeekday(time.Wednesday).String(), "WEDNESDAY")
		assert.Equal(t, NewWeekday(time.Thursday).String(), "THURSDAY")
		assert.Equal(t, NewWeekday(time.Friday).String(), "FRIDAY")
		assert.Equal(t, NewWeekday(time.Saturday).String(), "SATURDAY")
	})

	t.Run("GoWeekday", func(t *testing.T) {
		assert.Equal(t, NewWeekday(time.Sunday).GoWeekday(), time.Sunday)
		assert.Equal(t, NewWeekday(time.Monday).GoWeekday(), time.Monday)
		assert.Equal(t, NewWeekday(time.Tuesday).GoWeekday(), time.Tuesday)
		assert.Equal(t, NewWeekday(time.Wednesday).GoWeekday(), time.Wednesday)
		assert.Equal(t, NewWeekday(time.Thursday).GoWeekday(), time.Thursday)
		assert.Equal(t, NewWeekday(time.Friday).GoWeekday(), time.Friday)
		assert.Equal(t, NewWeekday(time.Saturday).GoWeekday(), time.Saturday)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		tests := []struct {
			weekday time.Weekday
			want    string
		}{
			{time.Sunday, `"SUNDAY"`},
			{time.Monday, `"MONDAY"`},
			{time.Tuesday, `"TUESDAY"`},
			{time.Wednesday, `"WEDNESDAY"`},
			{time.Thursday, `"THURSDAY"`},
			{time.Friday, `"FRIDAY"`},
			{time.Saturday, `"SATURDAY"`},
		}

		for _, tt := range tests {
			out, err := json.Marshal(NewWeekday(tt.weekday))
			assert.NoError(t, err)
			assert.Equal(t, string(out), tt.want)
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		tests := []struct {
			json string
			want Weekday
		}{
			{`"SUNDAY"`, NewWeekday(time.Sunday)},
			{`"MONDAY"`, NewWeekday(time.Monday)},
			{`"TUESDAY"`, NewWeekday(time.Tuesday)},
			{`"WEDNESDAY"`, NewWeekday(time.Wednesday)},
			{`"THURSDAY"`, NewWeekday(time.Thursday)},
			{`"FRIDAY"`, NewWeekday(time.Friday)},
			{`"SATURDAY"`, NewWeekday(time.Saturday)},
		}

		for _, tt := range tests {
			var wd Weekday
			err := json.Unmarshal([]byte(tt.json), &wd)
			assert.NoError(t, err)
			assert.Equal(t, wd, tt.want)
		}

		var wd Weekday
		err := json.Unmarshal([]byte(``), &wd)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`SUNDAY`), &wd)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`"SUNDAY1"`), &wd)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`0`), &wd)
		assert.Error(t, err)
	})
}

func TestTime(t *testing.T) {
	t.Run("NewTime", func(t *testing.T) {
		assert.Panic(t, func() { NewTime(0, 0, -1) })
		assert.Panic(t, func() { NewTime(0, 0, 61) })
		assert.Panic(t, func() { NewTime(0, -1, 0) })
		assert.Panic(t, func() { NewTime(0, 61, 0) })
		assert.Panic(t, func() { NewTime(-1, 0, 0) })
		assert.Panic(t, func() { NewTime(24, 0, 0) })
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewTime(0, 0, 0).String(), "00:00:00")
		assert.Equal(t, NewTime(1, 2, 3).String(), "01:02:03")
		assert.Equal(t, NewTime(23, 59, 59).String(), "23:59:59")
	})

	t.Run("Clock", func(t *testing.T) {
		hour, min, sec := NewTime(1, 2, 3).Clock()
		assert.Equal(t, hour, 1)
		assert.Equal(t, min, 2)
		assert.Equal(t, sec, 3)
	})

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, NewTime(0, 0, 0).IsZero())
		assert.False(t, NewTime(0, 0, 1).IsZero())
		assert.False(t, NewTime(0, 1, 0).IsZero())
		assert.False(t, NewTime(1, 0, 0).IsZero())
	})

	t.Run("Before", func(t *testing.T) {
		assert.True(t, NewTime(0, 0, 0).Before(NewTime(1, 0, 0)))
		assert.True(t, NewTime(0, 0, 0).Before(NewTime(0, 1, 0)))
		assert.True(t, NewTime(0, 0, 0).Before(NewTime(0, 0, 1)))
		assert.False(t, NewTime(1, 0, 0).Before(NewTime(0, 0, 0)))
		assert.False(t, NewTime(0, 1, 0).Before(NewTime(0, 0, 0)))
		assert.False(t, NewTime(0, 0, 1).Before(NewTime(0, 0, 0)))
		assert.False(t, NewTime(1, 1, 1).Before(NewTime(1, 1, 0)))
	})

	t.Run("After", func(t *testing.T) {
		assert.False(t, NewTime(0, 0, 0).After(NewTime(1, 0, 0)))
		assert.False(t, NewTime(0, 0, 0).After(NewTime(0, 1, 0)))
		assert.False(t, NewTime(0, 0, 0).After(NewTime(0, 0, 1)))
		assert.True(t, NewTime(1, 0, 0).After(NewTime(0, 0, 0)))
		assert.True(t, NewTime(0, 1, 0).After(NewTime(0, 0, 0)))
		assert.True(t, NewTime(0, 0, 1).After(NewTime(0, 0, 0)))
		assert.False(t, NewTime(1, 1, 0).After(NewTime(1, 1, 1)))
	})

	t.Run("Equal", func(t *testing.T) {
		assert.True(t, NewTime(0, 0, 0).Equal(NewTime(0, 0, 0)))
		assert.False(t, NewTime(0, 0, 0).Equal(NewTime(0, 0, 1)))
		assert.False(t, NewTime(0, 0, 0).Equal(NewTime(0, 1, 0)))
		assert.False(t, NewTime(0, 0, 0).Equal(NewTime(1, 0, 0)))
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		tm := NewTime(0, 0, 0)
		out, err := json.Marshal(tm)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"00:00:00"`)

		tm = NewTime(1, 2, 3)
		out, err = json.Marshal(tm)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"01:02:03"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var tm Time
		err := json.Unmarshal([]byte(`"00:00:00"`), &tm)
		assert.NoError(t, err)
		assert.Equal(t, tm, NewTime(0, 0, 0))

		err = json.Unmarshal([]byte(`"01:02:03"`), &tm)
		assert.NoError(t, err)
		assert.Equal(t, tm, NewTime(1, 2, 3))

		err = json.Unmarshal([]byte(`"01:02"`), &tm)
		assert.Error(t, err)
		err = json.Unmarshal([]byte(`"01:02:03:04"`), &tm)
		assert.Error(t, err)
	})
}

func TestDate(t *testing.T) {
	t.Run("NewDate", func(t *testing.T) {
		assert.Panic(t, func() { NewDate(time.Time{}.Year(), 0, 0) })
		assert.Panic(t, func() { NewDate(0, 0, 0) })
		assert.Panic(t, func() { NewDate(0, 0, -1) })
		assert.Panic(t, func() { NewDate(0, 0, 32) })
		assert.Panic(t, func() { NewDate(0, -1, 0) })
		assert.Panic(t, func() { NewDate(0, 13, 0) })
		assert.Panic(t, func() { NewDate(-1, 0, 0) })
		assert.Panic(t, func() { NewDate(1971, 1, 0) })
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewDate(2021, 1, 1).String(), "2021-01-01")
		assert.Equal(t, NewDate(2021, 12, 31).String(), "2021-12-31")
	})

	t.Run("Date", func(t *testing.T) {
		year, month, day := NewDate(2021, 1, 2).Date()
		assert.Equal(t, year, 2021)
		assert.Equal(t, month, 1)
		assert.Equal(t, day, 2)
	})

	t.Run("Before", func(t *testing.T) {
		assert.True(t, NewDate(2021, 1, 1).Before(NewDate(2021, 1, 2)))
		assert.True(t, NewDate(2021, 1, 1).Before(NewDate(2021, 2, 1)))
		assert.True(t, NewDate(2021, 1, 1).Before(NewDate(2022, 1, 1)))
		assert.False(t, NewDate(2021, 1, 2).Before(NewDate(2021, 1, 1)))
		assert.False(t, NewDate(2021, 2, 1).Before(NewDate(2021, 1, 1)))
		assert.False(t, NewDate(2022, 1, 1).Before(NewDate(2021, 1, 1)))
	})

	t.Run("After", func(t *testing.T) {
		assert.False(t, NewDate(2021, 1, 1).After(NewDate(2021, 1, 2)))
		assert.False(t, NewDate(2021, 1, 1).After(NewDate(2021, 2, 1)))
		assert.False(t, NewDate(2021, 1, 1).After(NewDate(2022, 1, 1)))
		assert.True(t, NewDate(2021, 1, 2).After(NewDate(2021, 1, 1)))
		assert.True(t, NewDate(2021, 2, 1).After(NewDate(2021, 1, 1)))
		assert.True(t, NewDate(2022, 1, 1).After(NewDate(2021, 1, 1)))
	})

	t.Run("Equal", func(t *testing.T) {
		assert.True(t, NewDate(2021, 1, 1).Equal(NewDate(2021, 1, 1)))
		assert.False(t, NewDate(2021, 1, 1).Equal(NewDate(2021, 1, 2)))
		assert.False(t, NewDate(2021, 1, 1).Equal(NewDate(2021, 2, 1)))
		assert.False(t, NewDate(2021, 1, 1).Equal(NewDate(2022, 1, 1)))
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		d := NewDate(2021, 1, 1)
		out, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"2021-01-01"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`"2021-01-01"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d, NewDate(2021, 1, 1))

		err = json.Unmarshal([]byte(`"2021-01"`), &d)
		assert.Error(t, err)
		err = json.Unmarshal([]byte(`"2021-01-01-01"`), &d)
		assert.Error(t, err)
	})
}

func TestDateTime(t *testing.T) {
	t.Run("NewDateTime", func(t *testing.T) {
		assert.Panic(t, func() { NewDateTime(0, 0, 0, 0, 0, -1) })
		assert.Panic(t, func() { NewDateTime(0, 0, 0, 0, 0, 61) })
		assert.Panic(t, func() { NewDateTime(0, 0, 0, 0, 61, 0) })
		assert.Panic(t, func() { NewDateTime(0, 0, 0, 24, 0, 0) })
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, NewDateTime(2021, 1, 1, 0, 0, 0).String(), "2021-01-01T00:00:00Z")
		assert.Equal(t, NewDateTime(2021, 12, 31, 23, 59, 59).String(), "2021-12-31T23:59:59Z")
	})

	t.Run("Date", func(t *testing.T) {
		year, month, day := NewDateTime(2021, 1, 2, 0, 0, 0).Date()
		assert.Equal(t, year, 2021)
		assert.Equal(t, month, 1)
		assert.Equal(t, day, 2)
	})

	t.Run("Clock", func(t *testing.T) {
		hour, min, sec := NewDateTime(2021, 1, 2, 3, 4, 5).Clock()
		assert.Equal(t, hour, 3)
		assert.Equal(t, min, 4)
		assert.Equal(t, sec, 5)
	})

	t.Run("GoTime", func(t *testing.T) {
		tm := NewDateTime(2021, 1, 2, 3, 4, 5)
		gt := tm.GoTime()
		assert.Equal(t, gt.Year(), 2021)
		assert.Equal(t, gt.Month(), time.January)
		assert.Equal(t, gt.Day(), 2)
		assert.Equal(t, gt.Hour(), 3)
		assert.Equal(t, gt.Minute(), 4)
		assert.Equal(t, gt.Second(), 5)
	})

	t.Run("Before", func(t *testing.T) {
		assert.True(t, NewDateTime(2021, 1, 1, 0, 0, 0).Before(NewDateTime(2021, 1, 1, 0, 0, 1)))
		assert.True(t, NewDateTime(2021, 1, 1, 0, 0, 0).Before(NewDateTime(2021, 1, 1, 0, 1, 0)))
		assert.True(t, NewDateTime(2021, 1, 1, 0, 0, 0).Before(NewDateTime(2021, 1, 1, 1, 0, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 1).Before(NewDateTime(2021, 1, 1, 0, 0, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 1, 0).Before(NewDateTime(2021, 1, 1, 0, 0, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 1, 0, 0).Before(NewDateTime(2021, 1, 1, 0, 0, 0)))
	})

	t.Run("After", func(t *testing.T) {
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).After(NewDateTime(2021, 1, 1, 0, 0, 1)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).After(NewDateTime(2021, 1, 1, 0, 1, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).After(NewDateTime(2021, 1, 1, 1, 0, 0)))
		assert.True(t, NewDateTime(2021, 1, 1, 0, 0, 1).After(NewDateTime(2021, 1, 1, 0, 0, 0)))
		assert.True(t, NewDateTime(2021, 1, 1, 0, 1, 0).After(NewDateTime(2021, 1, 1, 0, 0, 0)))
		assert.True(t, NewDateTime(2021, 1, 1, 1, 0, 0).After(NewDateTime(2021, 1, 1, 0, 0, 0)))
	})

	t.Run("Equal", func(t *testing.T) {
		assert.True(t, NewDateTime(2021, 1, 1, 0, 0, 0).Equal(NewDateTime(2021, 1, 1, 0, 0, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).Equal(NewDateTime(2021, 1, 1, 0, 0, 1)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).Equal(NewDateTime(2021, 1, 1, 0, 1, 0)))
		assert.False(t, NewDateTime(2021, 1, 1, 0, 0, 0).Equal(NewDateTime(2021, 1, 1, 1, 0, 0)))
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		dt := NewDateTime(2021, 1, 1, 0, 0, 0)
		out, err := json.Marshal(dt)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"2021-01-01T00:00:00Z"`)

		dt = NewDateTime(2021, 12, 31, 23, 59, 59)
		out, err = json.Marshal(dt)
		assert.NoError(t, err)
		assert.Equal(t, string(out), `"2021-12-31T23:59:59Z"`)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var dt DateTime
		err := json.Unmarshal([]byte(`"2021-01-01T00:00:00Z"`), &dt)
		assert.NoError(t, err)
		assert.Equal(t, dt, NewDateTime(2021, 1, 1, 0, 0, 0))

		err = json.Unmarshal([]byte(`"2021-01-01T00:00"`), &dt)
		assert.Error(t, err)
	})
}
