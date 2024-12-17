[![PkgGoDev](https://pkg.go.dev/badge/github.com/krhubert/timeapi)](https://pkg.go.dev/github.com/krhubert/timeapi)

# Time representation for API

This package provides a set of data types for time representation in API.

The problem with standard `time.Time` is that it is not clear what the time represents. Is it a date, a time, a date and time, a duration, an interval, etc.?

TimeApi address this issue by providing a set of data types that are more strict and clear about what they represent.

## Data types

1. Interval - represents a time interval between two points in time

2. Duration - represents the elapsed time

3. Timezone - represents a timezone with a name and offset

4. Weekday - represents a day of the week

5. Time - represents a time of the day

6. Date - represents a date

7. DateTime - represents a date and time

See the [documentation](https://pkg.go.dev/github.com/krhubert/timeapi) for more details.
