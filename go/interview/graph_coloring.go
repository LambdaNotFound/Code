package interview

/**
 * Greedy Graph Coloring
 *
 * Two beams that are too close in angle (< 10°)
 * on the same frequency will interfere
 *
 *   0° ── 3°
 *   |  ╲  |
 *   |   ╲ |
 *   8° ──(both connected to 0° and 3°)
 *   |
 *   15° ── 20° ── 25°
 *
 */

import (
	"fmt"
	"math"
	"sort"
)

type User struct {
	id    int
	angle float64
}

func assignColors(users []User, minSep float64, maxColors int) map[int]int {
	n := len(users)

	// Step 1: build conflict graph
	// conflicts[i] = list of user indices that conflict with user i
	conflicts := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := math.Abs(users[i].angle - users[j].angle)
			if diff < minSep {
				// these two beams would interfere — mark as conflicting
				conflicts[i] = append(conflicts[i], j)
				conflicts[j] = append(conflicts[j], i)
			}
		}
	}

	// Step 2: sort by most conflicts first (Welsh-Powell heuristic)
	// most constrained node is hardest to color → do it first
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return len(conflicts[order[i]]) > len(conflicts[order[j]])
	})

	// Step 3: greedy color assignment
	colors := make(map[int]int) // index → assigned color

	for _, i := range order {
		// find which colors are already used by this node's neighbors
		usedByNeighbors := map[int]bool{}
		for _, neighbor := range conflicts[i] {
			if c, assigned := colors[neighbor]; assigned {
				usedByNeighbors[c] = true
			}
		}

		// assign the lowest color not used by any neighbor
		for c := 1; c <= maxColors; c++ {
			if !usedByNeighbors[c] {
				colors[i] = c
				break
			}
		}

		// if no color found → this user can't be served
		if _, ok := colors[i]; !ok {
			fmt.Printf("  User %d (%.1f°) unserviceable — all colors conflict\n",
				users[i].id, users[i].angle)
		}
	}

	return colors
}

func testAssignColors() {
	users := []User{
		{1, 0.0},
		{2, 3.0},
		{3, 8.0},
		{4, 15.0},
		{5, 20.0},
		{6, 25.0},
	}

	colors := assignColors(users, 10.0, 4)

	fmt.Println("Assignments:")
	for i, u := range users {
		fmt.Printf("  User %d  angle %5.1f°  → color %d\n",
			u.id, u.angle, colors[i])
	}

	// verify no two conflicting users share a color
	fmt.Println("\nConflict check:")
	violations := 0
	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			diff := math.Abs(users[i].angle - users[j].angle)
			if diff < 10.0 && colors[i] == colors[j] {
				fmt.Printf("  VIOLATION: user %d and %d both color %d\n",
					users[i].id, users[j].id, colors[i])
				violations++
			}
		}
	}
	if violations == 0 {
		fmt.Println("  All constraints satisfied ✓")
	}
}
