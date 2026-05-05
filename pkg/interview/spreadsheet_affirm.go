package interview

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// -------- Typed Row Struct --------

type Row struct {
	Color  string
	Date   time.Time
	Number float64
}

type SpreadSheet struct {
	Rows []Row
}

// -------- Constructor --------

func NewSpreadSheet(filename string) (*SpreadSheet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Skip header line
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}

	var rows []Row
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // splits on any whitespace
		if len(fields) < 3 {
			continue
		}

		date, err := time.Parse("2006/01/02", fields[1])
		if err != nil {
			return nil, fmt.Errorf("invalid date %q: %w", fields[1], err)
		}

		number, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", fields[2], err)
		}

		rows = append(rows, Row{
			Color:  fields[0],
			Date:   date,
			Number: number,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &SpreadSheet{Rows: rows}, nil
}

// -------- Filter --------

type Condition struct {
	Column string
	Op     string
	Value  string
}

func (s *SpreadSheet) Filter(cond Condition) []Row {
	var result []Row
	for _, row := range s.Rows {
		if rowMatches(row, cond) {
			result = append(result, row)
		}
	}
	return result
}

func rowMatches(row Row, cond Condition) bool {
	switch cond.Column {
	case "color":
		return compareString(row.Color, cond.Op, cond.Value)

	case "date":
		filterDate, err := time.Parse("2006/01/02", cond.Value)
		if err != nil {
			return false
		}
		return compareDate(row.Date, cond.Op, filterDate)

	case "number":
		filterNum, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil {
			return false
		}
		return compareFloat(row.Number, cond.Op, filterNum)
	}
	return false
}

// -------- Comparators --------

func compareString(a, op, b string) bool {
	switch op {
	case "=":
		return a == b
	case "!=":
		return a != b
	}
	return false
}

func compareDate(a time.Time, op string, b time.Time) bool {
	switch op {
	case "=":
		return a.Equal(b)
	case "!=":
		return !a.Equal(b)
	case "<":
		return a.Before(b)
	case "<=":
		return a.Before(b) || a.Equal(b)
	case ">":
		return a.After(b)
	case ">=":
		return a.After(b) || a.Equal(b)
	}
	return false
}

func compareFloat(a float64, op string, b float64) bool {
	switch op {
	case "=":
		return a == b
	case "!=":
		return a != b
	case "<":
		return a < b
	case "<=":
		return a <= b
	case ">":
		return a > b
	case ">=":
		return a >= b
	}
	return false
}

// -------- Pretty Print --------

func printRows(rows []Row) {
	fmt.Printf("%-12s %-12s %-8s\n", "color", "date", "number")
	fmt.Println(strings.Repeat("-", 34))
	for _, row := range rows {
		fmt.Printf("%-12s %-12s %-8.0f\n",
			row.Color,
			row.Date.Format("2006/01/02"),
			row.Number,
		)
	}
}

// -------- Main --------

func main() {
	sheet, err := NewSpreadSheet("../fixtures/file.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== color = green ===")
	printRows(sheet.Filter(Condition{"color", "=", "green"}))

	fmt.Println("\n=== number < 5 ===")
	printRows(sheet.Filter(Condition{"number", "<", "5"}))

	fmt.Println("\n=== date > 2005/01/01 ===")
	printRows(sheet.Filter(Condition{"date", ">", "2005/01/01"}))
}
