package oodesign

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
 *
 * The Strategy Pattern
 *
 * The core idea: define an interface for a behavior, then swap in different implementations without changing the caller.
 *
 * In this code, the "behavior" is filtering a row:
 *
 *
 * Predicate interface        ← the strategy contract
 *    Match(row Row) bool
 *
 * EqualPredicate             ← concrete strategy A
 * NumericPredicate           ← concrete strategy B
 * AndPredicate               ← concrete strategy C (composes others)
 * SpreadSheet.Filter only knows about Predicate — it never asks "are you an EqualPredicate or a NumericPredicate?" It just calls Match. This means:
 *
 *
 * // These three lines look identical to Filter — it doesn't care which strategy is used
 * sheet.Filter(NewPredicate("color", "=", "green"))
 * sheet.Filter(NewPredicate("number", ">", "5"))
 * sheet.Filter(AndPredicate{...})
 *
 */

// -------- Row --------

type Row map[string]string

type SpreadSheet struct {
	Headers []string
	Rows    []Row
}

// -------- Predicate (Strategy + Composite) --------

type Predicate interface {
	Match(row Row) bool
}

type EqualPredicate struct{ Col, Val string }

func (p EqualPredicate) Match(row Row) bool { return row[p.Col] == p.Val }

type NotEqualPredicate struct{ Col, Val string }

func (p NotEqualPredicate) Match(row Row) bool { return row[p.Col] != p.Val }

// NumericPredicate handles ordering comparisons for numeric columns.
type NumericPredicate struct {
	Col string
	Op  string
	Val float64
}

func (p NumericPredicate) Match(row Row) bool {
	v, err := strconv.ParseFloat(row[p.Col], 64)
	if err != nil {
		return false
	}
	switch p.Op {
	case "<":
		return v < p.Val
	case "<=":
		return v <= p.Val
	case ">":
		return v > p.Val
	case ">=":
		return v >= p.Val
	}
	return false
}

// AndPredicate / OrPredicate compose any predicates without touching existing code.
type AndPredicate struct{ Preds []Predicate }

func (p AndPredicate) Match(row Row) bool {
	for _, pred := range p.Preds {
		if !pred.Match(row) {
			return false
		}
	}
	return true
}

type OrPredicate struct{ Preds []Predicate }

func (p OrPredicate) Match(row Row) bool {
	for _, pred := range p.Preds {
		if pred.Match(row) {
			return true
		}
	}
	return false
}

// NewPredicate is a factory that builds a predicate from a [col, op, val] triple.
func NewPredicate(col, op, val string) Predicate {
	switch op {
	case "=":
		return EqualPredicate{col, val}
	case "!=":
		return NotEqualPredicate{col, val}
	case "<", "<=", ">", ">=":
		v, _ := strconv.ParseFloat(val, 64)
		return NumericPredicate{col, op, v}
	}
	return nil
}

// -------- SpreadSheet --------

func NewSpreadSheet(filename string) (*SpreadSheet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	headers := strings.Fields(scanner.Text())

	var rows []Row
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != len(headers) {
			continue
		}
		row := Row{}
		for i, h := range headers {
			row[h] = fields[i]
		}
		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &SpreadSheet{Headers: headers, Rows: rows}, nil
}

func (s *SpreadSheet) Filter(pred Predicate) []Row {
	var result []Row
	for _, row := range s.Rows {
		if pred.Match(row) {
			result = append(result, row)
		}
	}
	return result
}

// -------- Pretty Print --------

func printRows(rows []Row, headers []string) {
	for _, h := range headers {
		fmt.Printf("%-14s", h)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 14*len(headers)))
	for _, row := range rows {
		for _, h := range headers {
			fmt.Printf("%-14s", row[h])
		}
		fmt.Println()
	}
}

// -------- Demo --------

func main() {
	sheet, err := NewSpreadSheet("../fixtures/file.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== color = green ===")
	printRows(sheet.Filter(NewPredicate("color", "=", "green")), sheet.Headers)

	fmt.Println("\n=== number < 5 ===")
	printRows(sheet.Filter(NewPredicate("number", "<", "5")), sheet.Headers)

	fmt.Println("\n=== color = green AND number > 5 ===")
	printRows(sheet.Filter(AndPredicate{[]Predicate{
		NewPredicate("color", "=", "green"),
		NewPredicate("number", ">", "5"),
	}}), sheet.Headers)

	fmt.Println("\n=== color = green OR color = white ===")
	printRows(sheet.Filter(OrPredicate{[]Predicate{
		NewPredicate("color", "=", "green"),
		NewPredicate("color", "=", "white"),
	}}), sheet.Headers)
}
