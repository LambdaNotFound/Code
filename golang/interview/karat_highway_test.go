package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func logFileFromLines(lines []string) *LogFile {
	lf := &LogFile{}
	for _, line := range lines {
		lf.LogEntries = append(lf.LogEntries, NewLogEntry(line))
	}
	return lf
}

func Test_NewLogEntry(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		timestamp float64
		plate     string
		location  int
		direction string
		boothType string
	}{
		{
			name:      "mainroad_east",
			line:      "44776.619 KTB918 310E MAINROAD",
			timestamp: 44776.619,
			plate:     "KTB918",
			location:  310,
			direction: "EAST",
			boothType: "MAINROAD",
		},
		{
			name:      "entry_west",
			line:      "52160.132 ABC123 400W ENTRY",
			timestamp: 52160.132,
			plate:     "ABC123",
			location:  400,
			direction: "WEST",
			boothType: "ENTRY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogEntry(tt.line)
			assert.Equal(t, tt.timestamp, got.Timestamp)
			assert.Equal(t, tt.plate, got.LicensePlate)
			assert.Equal(t, tt.location, got.Location)
			assert.Equal(t, tt.direction, got.Direction)
			assert.Equal(t, tt.boothType, got.BoothType)
		})
	}
}

func Test_CountJourneys(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			name: "single_direct_journey",
			lines: []string{
				"1.000 AAA111 100E ENTRY",
				"2.000 AAA111 110E EXIT",
			},
			want: 1,
		},
		{
			name: "three_journeys_interleaved",
			// JOX304: 1 journey; THX138: 2 journeys
			lines: []string{
				"90750.191 JOX304 250E ENTRY",
				"91081.684 JOX304 260E MAINROAD",
				"91082.101 THX138 110E ENTRY",
				"91483.251 JOX304 270E MAINROAD",
				"91873.920 THX138 120E MAINROAD",
				"91874.493 JOX304 280E EXIT",
				"91982.102 THX138 290E EXIT",
				"92301.302 THX138 300E ENTRY",
				"92371.302 THX138 310E EXIT",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, logFileFromLines(tt.lines).CountJourneys())
		})
	}
}

func Test_CatchSpeeders(t *testing.T) {
	// JOX304: safe (max 109 km/h) — control group
	// TST002: one segment at ~131 km/h → 1 ticket (≥130 rule)
	// TST003 journey 1: two segments at 125 km/h → 1 ticket (≥120 in 2 segments rule)
	// TST003 journey 2: one segment at ~130 km/h → 1 ticket (≥130 rule)
	lines := []string{
		"0.000 TST003 100E ENTRY",
		"288.000 TST003 110E MAINROAD",  // speed=125 km/h, highCount=1
		"576.000 TST003 120E EXIT",      // speed=125 km/h, highCount=2 → flagged; journey reset
		"1000.000 TST003 100E ENTRY",
		"1276.000 TST003 110E EXIT",     // speed≈130.4 km/h → flagged; journey reset
		"90750.191 JOX304 250E ENTRY",
		"91081.684 JOX304 260E MAINROAD", // speed≈109 km/h
		"91483.251 JOX304 270E MAINROAD", // speed≈90 km/h
		"91874.493 JOX304 280E EXIT",     // speed≈92 km/h
		"93005.405 TST002 270W ENTRY",
		"93280.609 TST002 260W EXIT",    // speed≈131 km/h → flagged
	}

	t.Run("two_plates_three_tickets", func(t *testing.T) {
		result := logFileFromLines(lines).CatchSpeeders()
		counts := map[string]int{}
		for _, plate := range result {
			counts[plate]++
		}
		assert.Equal(t, 1, counts["TST002"])
		assert.Equal(t, 2, counts["TST003"])
		assert.Equal(t, 0, counts["JOX304"])
		assert.Len(t, counts, 2) // only TST002 and TST003
	})

	t.Run("empty_log", func(t *testing.T) {
		assert.Empty(t, logFileFromLines(nil).CatchSpeeders())
	})

	t.Run("single_booth_no_segment", func(t *testing.T) {
		// Only an ENTRY — no segment to measure, no ticket
		result := logFileFromLines([]string{
			"1.000 XYZ999 100E ENTRY",
		}).CatchSpeeders()
		assert.Empty(t, result)
	})

	t.Run("exactly_120_two_segments", func(t *testing.T) {
		// 36000/300 = 120 km/h exactly in both segments → highCount=2 → speeder
		result := logFileFromLines([]string{
			"0.000 PLT001 100E ENTRY",
			"300.000 PLT001 110E MAINROAD",
			"600.000 PLT001 120E EXIT",
		}).CatchSpeeders()
		assert.Equal(t, []string{"PLT001"}, result)
	})

	t.Run("exactly_120_one_segment_not_speeder", func(t *testing.T) {
		// Only one segment at exactly 120 km/h — below both thresholds
		result := logFileFromLines([]string{
			"0.000 PLT002 100E ENTRY",
			"300.000 PLT002 110E EXIT",
		}).CatchSpeeders()
		assert.Empty(t, result)
	})

	t.Run("flagged_once_per_journey", func(t *testing.T) {
		// Three ≥130 segments in one journey — plate should appear only once
		result := logFileFromLines([]string{
			"0.000 PLT003 100E ENTRY",
			"270.000 PLT003 110E MAINROAD",  // speed≈133 km/h
			"540.000 PLT003 120E MAINROAD",  // speed≈133 km/h
			"810.000 PLT003 130E EXIT",      // speed≈133 km/h
		}).CatchSpeeders()
		assert.Equal(t, []string{"PLT003"}, result)
	})
}
