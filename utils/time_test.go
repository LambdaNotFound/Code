package utils

import "time"

// time.RFC3339     // "2006-01-02T15:04:05Z07:00"
// time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"
// time.DateTime    // "2006-01-02 15:04:05"  (Go 1.20+)
// time.DateOnly    // "2006-01-02"           (Go 1.20+)
// time.TimeOnly    // "15:04:05"             (Go 1.20+)

/*
Why it's easy to remember
The reference values count up: 1 2 3 4 5 6 7
month   = 1   (January)
day     = 2
hour    = 3   (15 in 24h)
minute  = 4   (04)
second  = 5   (05)
year    = 6   (2006)
weekday = 7   (Monday — 7th day if Sunday=0, or just remember Mon)
*/

func parse_date(str string) time.Time {
	t, _ := time.Parse("2006/01/02", "2001/02/23") // slash separated
	t, _ = time.Parse("2006/01/02", str)

	t, _ = time.Parse("2006-01-02", "2001-02-23")   // dash separated
	t, _ = time.Parse("01-02-2006", "02-23-2001")   // mm-dd-yyyy
	t, _ = time.Parse("2006 01 02", "2001 02 23")   // space separated
	t, _ = time.Parse("Jan 02 2006", "Feb 23 2001") // named month

	return t
}
