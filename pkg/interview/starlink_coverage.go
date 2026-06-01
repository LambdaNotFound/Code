package interview

import (
	"fmt"
	"sort"
)

type Station struct {
	pos, reach float64
	id         int
}

func (s Station) left() float64  { return s.pos - s.reach }
func (s Station) right() float64 { return s.pos + s.reach }

// minStationsCover returns the minimum stations to cover [passStart, passEnd]
// with no gaps, and which ones. Returns (-1, nil) if full coverage is impossible.
func minStationsCover(passStart, passEnd float64, stations []Station) (int, []Station) {
	// sort by left edge of coverage
	sort.Slice(stations, func(i, j int) bool { return stations[i].left() < stations[j].left() })

	var chosen []Station
	frontier := passStart // everything up to `frontier` is already covered
	i, n := 0, len(stations)

	for frontier < passEnd {
		// among all stations reaching the current frontier, take the one
		// that extends coverage the furthest right
		bestRight := frontier
		bestIdx := -1
		for i < n && stations[i].left() <= frontier {
			if r := stations[i].right(); r > bestRight {
				bestRight, bestIdx = r, i
			}
			i++
		}
		if bestIdx == -1 {
			return -1, nil // gap: nothing covers the point at `frontier`
		}
		chosen = append(chosen, stations[bestIdx])
		frontier = bestRight
	}
	return len(chosen), chosen
}

func minStationsCoverTest() {
	passStart, passEnd := 0.0, 1000.0

	stations := []Station{
		{id: 1, pos: 100, reach: 120},
		{id: 2, pos: 300, reach: 150},
		{id: 3, pos: 250, reach: 80}, // redundant
		{id: 4, pos: 550, reach: 180},
		{id: 5, pos: 800, reach: 220},
		{id: 6, pos: 700, reach: 50}, // redundant
	}

	count, chosen := minStationsCover(passStart, passEnd, stations)
	if count == -1 {
		fmt.Println("Coverage impossible — there is a gap in the pass.")
		return
	}
	fmt.Printf("Minimum stations needed: %d\n", count)
	fmt.Println("Chosen stations:")
	for _, s := range chosen {
		fmt.Printf("  Station %d: pos %.0f, reach ±%.0f -> covers [%.0f, %.0f]\n",
			s.id, s.pos, s.reach, s.left(), s.right())
	}
}
