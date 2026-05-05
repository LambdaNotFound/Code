package interview

import (
	"os"
	"testing"
	"time"
)

// helper to create a temp CSV file for testing
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "spreadsheet_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleCSV = `color	date	number
green	2001/02/23	8
purple	2006/05/11	1
white	2019/02/17	200
green	2002/02/23	28`

func mustDate(s string) time.Time {
	t, _ := time.Parse("2006/01/02", s)
	return t
}

func TestFilter(t *testing.T) {
	filename := createTempFile(t, sampleCSV)
	defer os.Remove(filename)

	sheet, err := NewSpreadSheet(filename)
	if err != nil {
		t.Fatalf("failed to load spreadsheet: %v", err)
	}

	tests := []struct {
		name      string
		condition Condition
		want      []Row
	}{
		{
			name:      "color equals green",
			condition: Condition{"color", "=", "green"},
			want: []Row{
				{Color: "green", Date: mustDate("2001/02/23"), Number: 8},
				{Color: "green", Date: mustDate("2002/02/23"), Number: 28},
			},
		},
		{
			name:      "color not equals green",
			condition: Condition{"color", "!=", "green"},
			want: []Row{
				{Color: "purple", Date: mustDate("2006/05/11"), Number: 1},
				{Color: "white", Date: mustDate("2019/02/17"), Number: 200},
			},
		},
		{
			name:      "number less than 5",
			condition: Condition{"number", "<", "5"},
			want: []Row{
				{Color: "purple", Date: mustDate("2006/05/11"), Number: 1},
			},
		},
		{
			name:      "number greater than or equal to 8",
			condition: Condition{"number", ">=", "8"},
			want: []Row{
				{Color: "green", Date: mustDate("2001/02/23"), Number: 8},
				{Color: "white", Date: mustDate("2019/02/17"), Number: 200},
				{Color: "green", Date: mustDate("2002/02/23"), Number: 28},
			},
		},
		{
			name:      "date after 2005",
			condition: Condition{"date", ">", "2005/01/01"},
			want: []Row{
				{Color: "purple", Date: mustDate("2006/05/11"), Number: 1},
				{Color: "white", Date: mustDate("2019/02/17"), Number: 200},
			},
		},
		{
			name:      "date equals specific date",
			condition: Condition{"date", "=", "2019/02/17"},
			want: []Row{
				{Color: "white", Date: mustDate("2019/02/17"), Number: 200},
			},
		},
		{
			name:      "no match returns empty",
			condition: Condition{"color", "=", "blue"},
			want:      nil,
		},
		{
			name:      "unknown column returns empty",
			condition: Condition{"unknown", "=", "foo"},
			want:      nil,
		},
		{
			name:      "invalid date value returns empty",
			condition: Condition{"date", ">", "not-a-date"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sheet.Filter(tt.condition)
			if !rowsEqual(got, tt.want) {
				t.Errorf("\ncondition: %+v\ngot:  %+v\nwant: %+v", tt.condition, got, tt.want)
			}
		})
	}
}

// rowsEqual compares two Row slices for deep equality
func rowsEqual(a, b []Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Color != b[i].Color ||
			!a[i].Date.Equal(b[i].Date) ||
			a[i].Number != b[i].Number {
			return false
		}
	}
	return true
}
