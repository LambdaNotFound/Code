package interview

import (
	"os"
	"testing"
)

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

func TestFilter(t *testing.T) {
	filename := createTempFile(t, sampleCSV)
	defer os.Remove(filename)

	sheet, err := NewSpreadSheet(filename)
	if err != nil {
		t.Fatalf("failed to load spreadsheet: %v", err)
	}

	tests := []struct {
		name string
		pred Predicate
		want []Row
	}{
		{
			name: "color equals green",
			pred: NewPredicate("color", "=", "green"),
			want: []Row{
				{"color": "green", "date": "2001/02/23", "number": "8"},
				{"color": "green", "date": "2002/02/23", "number": "28"},
			},
		},
		{
			name: "color not equals green",
			pred: NewPredicate("color", "!=", "green"),
			want: []Row{
				{"color": "purple", "date": "2006/05/11", "number": "1"},
				{"color": "white", "date": "2019/02/17", "number": "200"},
			},
		},
		{
			name: "number less than 5",
			pred: NewPredicate("number", "<", "5"),
			want: []Row{
				{"color": "purple", "date": "2006/05/11", "number": "1"},
			},
		},
		{
			name: "number greater than or equal to 8",
			pred: NewPredicate("number", ">=", "8"),
			want: []Row{
				{"color": "green", "date": "2001/02/23", "number": "8"},
				{"color": "white", "date": "2019/02/17", "number": "200"},
				{"color": "green", "date": "2002/02/23", "number": "28"},
			},
		},
		{
			name: "date equals specific date",
			pred: NewPredicate("date", "=", "2019/02/17"),
			want: []Row{
				{"color": "white", "date": "2019/02/17", "number": "200"},
			},
		},
		{
			name: "no match returns nil",
			pred: NewPredicate("color", "=", "blue"),
			want: nil,
		},
		{
			name: "unknown column returns nil",
			pred: NewPredicate("unknown", "=", "foo"),
			want: nil,
		},
		{
			name: "AND: green and number > 5",
			pred: AndPredicate{[]Predicate{
				NewPredicate("color", "=", "green"),
				NewPredicate("number", ">", "5"),
			}},
			want: []Row{
				{"color": "green", "date": "2001/02/23", "number": "8"},
				{"color": "green", "date": "2002/02/23", "number": "28"},
			},
		},
		{
			name: "OR: green or white",
			pred: OrPredicate{[]Predicate{
				NewPredicate("color", "=", "green"),
				NewPredicate("color", "=", "white"),
			}},
			want: []Row{
				{"color": "green", "date": "2001/02/23", "number": "8"},
				{"color": "white", "date": "2019/02/17", "number": "200"},
				{"color": "green", "date": "2002/02/23", "number": "28"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sheet.Filter(tt.pred)
			if !rowsEqual(got, tt.want) {
				t.Errorf("\npred: %+v\ngot:  %+v\nwant: %+v", tt.pred, got, tt.want)
			}
		})
	}
}

func rowsEqual(a, b []Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for k, v := range a[i] {
			if b[i][k] != v {
				return false
			}
		}
	}
	return true
}
