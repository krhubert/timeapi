// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package timeapi

import (
	"errors"
	"math"
	"strconv"
	"time"
)

var durationUnitMap = map[string]uint64{
	"h": uint64(time.Hour),
	"m": uint64(time.Minute),
	"s": uint64(time.Second),
}

var durationUnitRank = map[string]int{
	"h": 0,
	"m": 1,
	"s": 2,
}

// ParseDuration is a modified version of time.ParseDuration
// A duration string is a possibly signed sequence of
// decimal numbers, each with a unit suffix,
// such as "300s", "-1h" or "2h45m".
// Valid time units are "s", "m", "h".
// This implementation disallows:
//   - of "ns", "us", "ms" and "d" units.
//   - of fractions of a unit.
//   - repeating a unit.
//   - units out of order.
func ParseDuration(s string) (Duration, error) {
	// [-+]?([0-9]*[a-z]+)+
	orig := s
	var dtmp uint64
	var d Duration

	seen := map[string]bool{
		"s": false,
		"m": false,
		"h": false,
	}
	maxRank := -1
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}

	if s == "" {
		return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
	}

	d.neg = 1
	if neg {
		d.neg = -1
	}

	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return d, nil
	}

	for s != "" {
		var v uint64
		var err error

		// The next character must be [0-9.]
		if !('0' <= s[0] && s[0] <= '9') {
			return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
		}
		// Consume [0-9]*
		v, s, err = leadingInt(s)
		if err != nil {
			return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return d, errors.New("timeapi: missing unit in duration " + strconv.Quote(orig))
		}

		u := s[:i]
		s = s[i:]
		unit, ok := durationUnitMap[u]
		if !ok {
			return d, errors.New("timeapi: unknown unit " + strconv.Quote(u) + " in duration " + strconv.Quote(orig))
		}
		if seen[u] {
			return d, errors.New("timeapi: unit " + strconv.Quote(u) + " repeated in duration " + strconv.Quote(orig))
		}

		// make sure unit is in order h, m, s
		seenRank := durationUnitRank[u]
		if seenRank < maxRank {
			return d, errors.New("timeapi: unit " + strconv.Quote(u) + " must be in the order of h, m, s in duration " + strconv.Quote(orig))
		}

		maxRank = max(seenRank, maxRank)
		seen[u] = true

		if v > 1<<63/unit {
			// overflow
			return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
		}

		switch u {
		case "h":
			d.hour = int(v)
		case "m":
			d.minute = int(v)
		case "s":
			d.second = int(v)
		}
		v *= unit

		dtmp += v
		if dtmp > 1<<63 {
			return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
		}
	}
	if neg {
		return d, nil
	}
	if dtmp > 1<<63-1 {
		return d, errors.New("timeapi: invalid duration " + strconv.Quote(orig))
	}
	return d, nil
}

var intervalUnitRank = map[string]int{
	"y":  0,
	"mo": 1,
	"d":  2,
	"h":  3,
	"m":  4,
	"s":  5,
}

// ParseInterval parses a string and returns an Interval.
// A interval string is a sequence of decimal numbers,
// each with a unit suffix, such as "1y", "1mo" or "2h45m".
// Valid units are "y", "mo", "d", "h", "m", "s".
func ParseInterval(s string) (Interval, error) {
	// ([0-9]*[a-z]+)+
	orig := s
	seen := map[string]bool{
		"y":  false,
		"mo": false,
		"d":  false,
		"h":  false,
		"m":  false,
		"s":  false,
	}

	maxRank := -1
	var ivl Interval

	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return ivl, nil
	}
	if s == "" {
		return ivl, errors.New("timeapi: invalid interval " + strconv.Quote(orig))
	}

	for s != "" {
		var v uint64
		var err error

		// The next character must be [0-9]
		if !('0' <= s[0] && s[0] <= '9') {
			return ivl, errors.New("timeapi: invalid interval " + strconv.Quote(orig))
		}
		// Consume [0-9]*
		v, s, err = leadingInt(s)
		if err != nil {
			return ivl, errors.New("timeapi: invalid interval " + strconv.Quote(orig))
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if '0' <= c && c <= '9' {
				break
			}
		}

		if i == 0 {
			return ivl, errors.New("timeapi: missing unit in interval " + strconv.Quote(orig))
		}

		u := s[:i]
		s = s[i:]
		if _, ok := seen[u]; !ok {
			return ivl, errors.New("timeapi: unknown unit " + strconv.Quote(u) + " in interval " + strconv.Quote(orig))
		}

		if seen[u] {
			return ivl, errors.New("timeapi: unit " + strconv.Quote(u) + " repeated in interval " + strconv.Quote(orig))
		}

		// make sure unit is in order
		seenRank := intervalUnitRank[u]
		if seenRank < maxRank {
			return ivl, errors.New("timeapi: unit " + strconv.Quote(u) + " must be in the order of y, mo, d, h, m, s in interval " + strconv.Quote(orig))
		}

		maxRank = max(seenRank, maxRank)
		seen[u] = true

		if v > math.MaxInt {
			// overflow
			return ivl, errors.New("timeapi: invalid interval " + strconv.Quote(orig))
		}

		switch u {
		case "y":
			ivl.year = int(v)
		case "mo":
			ivl.month = int(v)
		case "d":
			ivl.day = int(v)
		case "h":
			ivl.hour = int(v)
		case "m":
			ivl.minute = int(v)
		case "s":
			ivl.second = int(v)
		}
	}
	return ivl, nil
}

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errors.New("timeapi: bad [0-9]*")
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errors.New("timeapi: bad [0-9]*")
		}
	}
	return x, s[i:], nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
