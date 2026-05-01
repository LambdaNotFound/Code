package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── Mock ──────────────────────────────────────────────────────────────────────

// mockFile replaces the global read4 with a closure that simulates reading
// sequentially from src, returning at most 4 bytes per call — exactly as the
// LeetCode judge API behaves. Returns a restore func that resets read4 to a
// no-op so tests don't bleed state into each other.
func mockFile(src string) (restore func()) {
	pos := 0
	read4 = func(buf4 []byte) int {
		// copy up to 4 bytes from src starting at pos
		n := copy(buf4[:4], src[pos:])
		pos += n
		return n
	}
	return func() { read4 = func(buf4 []byte) int { return 0 } }
}

// ── ReadBuffer mock ───────────────────────────────────────────────────────────

// mockPaginator replaces the global fetchPage with a closure that simulates a
// paginated integer source backed by pages (a slice of slices). Page indices are
// used directly as cursors, matching ReadBuffer.cursor semantics. Returns a
// restore func that resets fetchPage to a no-op so tests don't bleed state.
func mockPaginator(pages [][]int) (restore func()) {
	fetchPage = func(page int) Result {
		if page >= len(pages) {
			return Result{nextPage: nil, results: nil}
		}
		var next *int
		if page+1 < len(pages) {
			n := page + 1
			next = &n
		}
		return Result{nextPage: next, results: pages[page]}
	}
	return func() { fetchPage = func(page int) Result { return Result{} } }
}

// ── ReadBuffer.fetch ──────────────────────────────────────────────────────────

func Test_ReadBuffer_fetch(t *testing.T) {
	testCases := []struct {
		name     string
		pages    [][]int // simulated paginated source
		requests []int   // successive fetch(n) calls
		want     [][]int // expected result per call
	}{
		{
			// Single page, fetch less than available.
			name:     "single_page_fetch_partial",
			pages:    [][]int{{1, 2, 3, 4, 5}},
			requests: []int{3},
			want:     [][]int{{1, 2, 3}},
		},
		{
			// Single page, fetch exactly all available.
			name:     "single_page_fetch_exact",
			pages:    [][]int{{10, 20, 30}},
			requests: []int{3},
			want:     [][]int{{10, 20, 30}},
		},
		{
			// Single page, request more than available → returns what exists.
			name:     "single_page_fetch_over",
			pages:    [][]int{{7, 8}},
			requests: []int{10},
			want:     [][]int{{7, 8}},
		},
		{
			// Two pages, single fetch spanning both.
			name:     "two_pages_single_fetch",
			pages:    [][]int{{1, 2, 3}, {4, 5, 6}},
			requests: []int{6},
			want:     [][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			// Two pages, fetch stops mid-first-page; second fetch drains rest + page2.
			name:     "two_pages_two_fetches",
			pages:    [][]int{{1, 2, 3}, {4, 5, 6}},
			requests: []int{2, 4},
			want:     [][]int{{1, 2}, {3, 4, 5, 6}},
		},
		{
			// Three pages, multiple fetches crossing all boundaries.
			name:     "three_pages_multiple_fetches",
			pages:    [][]int{{10, 20}, {30, 40}, {50, 60}},
			requests: []int{1, 3, 2},
			want:     [][]int{{10}, {20, 30, 40}, {50, 60}},
		},
		{
			// Empty source: every fetch returns nothing.
			name:     "empty_source",
			pages:    [][]int{},
			requests: []int{5, 5},
			want:     [][]int{{}, {}},
		},
		{
			// Fetch after exhaustion returns empty.
			name:     "fetch_after_exhaustion",
			pages:    [][]int{{1, 2}},
			requests: []int{2, 3},
			want:     [][]int{{1, 2}, {}},
		},
		{
			// One item at a time across multiple pages.
			name:     "one_item_at_a_time",
			pages:    [][]int{{1, 2}, {3}},
			requests: []int{1, 1, 1, 1},
			want:     [][]int{{1}, {2}, {3}, {}},
		},
		{
			// Page with zero items followed by a page with items.
			// The empty page causes fetchedOffset==0, so the loop re-enters the
			// refill branch immediately and fetches the next page.
			name:     "empty_page_then_data",
			pages:    [][]int{{}, {9, 8, 7}},
			requests: []int{2},
			want:     [][]int{{9, 8}},
		},
		{
			// Large request across many small pages.
			name:     "many_small_pages",
			pages:    [][]int{{1}, {2}, {3}, {4}, {5}},
			requests: []int{5},
			want:     [][]int{{1, 2, 3, 4, 5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			restore := mockPaginator(tc.pages)
			defer restore()

			rb := &ReadBuffer{}
			for i, n := range tc.requests {
				got := rb.fetch(n)
				assert.Equal(t, tc.want[i], got, "fetch[%d](n=%d)", i, n)
			}
		})
	}
}

// ── 157: Read N Characters Given Read4 (single call) ─────────────────────────

func Test_read157(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		n       int
		wantN   int
		wantBuf string
	}{
		{
			name: "read_exact_file_length",
			// file is 3 chars, read 3 — consumes one read4 call returning 3
			file: "abc", n: 3,
			wantN: 3, wantBuf: "abc",
		},
		{
			name: "read_less_than_file_length",
			// file is 5 chars, read 3 — stops after satisfying n
			file: "abcde", n: 3,
			wantN: 3, wantBuf: "abc",
		},
		{
			name: "read_more_than_file_length",
			// file is 3 chars, read 6 — hits EOF, returns 3
			file: "abc", n: 6,
			wantN: 3, wantBuf: "abc",
		},
		{
			name: "read_exactly_4_chars_one_read4_call",
			// file is 4 chars — one full read4 call, n satisfied exactly
			file: "abcd", n: 4,
			wantN: 4, wantBuf: "abcd",
		},
		{
			name: "read_across_two_read4_calls",
			// file is 8 chars, read 8 — two full read4 calls
			file: "abcdefgh", n: 8,
			wantN: 8, wantBuf: "abcdefgh",
		},
		{
			name: "read_n_larger_spans_multiple_read4_calls",
			// file is 12 chars, read 10 — three read4 calls (4+4+4), stop after 10
			file: "abcdefghijkl", n: 10,
			wantN: 10, wantBuf: "abcdefghij",
		},
		{
			name: "empty_file_returns_zero",
			// file is empty — read4 immediately returns 0 (EOF)
			file: "", n: 5,
			wantN: 0, wantBuf: "",
		},
		{
			name: "single_char_file",
			file: "x", n: 1,
			wantN: 1, wantBuf: "x",
		},
		{
			name: "read_1_from_large_file",
			// only first char is needed, subsequent read4 data is discarded
			file: "abcdefghij", n: 1,
			wantN: 1, wantBuf: "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			restore := mockFile(tc.file)
			defer restore()

			buf := make([]byte, tc.n+1) // +1 to detect over-writes
			got := read(buf, tc.n)

			assert.Equal(t, tc.wantN, got, "chars returned")
			assert.Equal(t, tc.wantBuf, string(buf[:got]), "content")
		})
	}
}

// ── 158: Read N Characters Given read4 II — Call Multiple Times ──────────────

func Test_Read4_read158(t *testing.T) {
	type call struct {
		n       int
		wantN   int
		wantBuf string
	}

	testCases := []struct {
		name  string
		file  string
		calls []call
	}{
		{
			// LeetCode Example 1
			// file="abc", calls: read(1)→"a", read(2)→"bc", read(4)→""
			// The first read4 fetches all 3 chars into buf4.
			// read(1) drains 1 char, leaving "bc" in buf4 for the next call.
			// read(2) drains the remaining 2 chars.
			// read(4) triggers another read4 which returns 0 (EOF).
			name: "leetcode_example_1",
			file: "abc",
			calls: []call{
				{n: 1, wantN: 1, wantBuf: "a"},
				{n: 2, wantN: 2, wantBuf: "bc"},
				{n: 4, wantN: 0, wantBuf: ""},
			},
		},
		{
			// LeetCode Example 2
			// file="abc", calls: read(4)→"abc", read(1)→""
			name: "leetcode_example_2",
			file: "abc",
			calls: []call{
				{n: 4, wantN: 3, wantBuf: "abc"},
				{n: 1, wantN: 0, wantBuf: ""},
			},
		},
		{
			// State preserved across calls spanning multiple read4 chunks.
			// file="abcdefgh" (8 chars, two read4 calls)
			// read(3) → drains "abc" from chunk 1, leaves "d" in buf4
			// read(5) → drains "d", fetches chunk 2 "efgh", drains "efgh" = 5 chars
			name: "state_persists_across_chunk_boundary",
			file: "abcdefgh",
			calls: []call{
				{n: 3, wantN: 3, wantBuf: "abc"},
				{n: 5, wantN: 5, wantBuf: "defgh"},
			},
		},
		{
			// Internal buf4 leftovers span two read calls.
			// file="abcde" (5 chars → read4 yields "abcd" then "e")
			// read(2) → consumes "ab", "cd" remains in buf4
			// read(3) → drains "cd" from buf4, fetches "e" from read4, returns "cde"
			name: "leftover_buf4_consumed_on_next_call",
			file: "abcde",
			calls: []call{
				{n: 2, wantN: 2, wantBuf: "ab"},
				{n: 3, wantN: 3, wantBuf: "cde"},
			},
		},
		{
			// Many small reads across a file larger than 4 chars.
			// Exercises the buf4 refill path multiple times.
			name: "many_small_reads",
			file: "abcdefghij", // 10 chars
			calls: []call{
				{n: 1, wantN: 1, wantBuf: "a"},    // buf4="abcd", idx→1
				{n: 1, wantN: 1, wantBuf: "b"},    // buf4 leftover, idx→2
				{n: 2, wantN: 2, wantBuf: "cd"},   // buf4 leftover drains exactly
				{n: 4, wantN: 4, wantBuf: "efgh"}, // new read4 chunk
				{n: 4, wantN: 2, wantBuf: "ij"},   // partial last chunk
				{n: 1, wantN: 0, wantBuf: ""},     // EOF
			},
		},
		{
			// Read exactly on chunk boundaries.
			name: "reads_aligned_to_chunk_boundaries",
			file: "abcdefgh",
			calls: []call{
				{n: 4, wantN: 4, wantBuf: "abcd"},
				{n: 4, wantN: 4, wantBuf: "efgh"},
				{n: 4, wantN: 0, wantBuf: ""},
			},
		},
		{
			// Request more than the file has on every call.
			name: "read_exceeds_file_every_call",
			file: "hi",
			calls: []call{
				{n: 100, wantN: 2, wantBuf: "hi"},
				{n: 100, wantN: 0, wantBuf: ""},
			},
		},
		{
			// Empty file — every read returns 0 immediately.
			name: "empty_file",
			file: "",
			calls: []call{
				{n: 5, wantN: 0, wantBuf: ""},
				{n: 5, wantN: 0, wantBuf: ""},
			},
		},
		{
			// Single character file with multiple read attempts.
			name: "single_char_file_multi_reads",
			file: "z",
			calls: []call{
				{n: 1, wantN: 1, wantBuf: "z"},
				{n: 1, wantN: 0, wantBuf: ""},
			},
		},
		{
			// Read 1 char at a time through a 5-char file.
			// buf4 is filled once ("abcd") and drained one char per call,
			// then refilled for the 5th char.
			name: "one_char_at_a_time",
			file: "abcde",
			calls: []call{
				{n: 1, wantN: 1, wantBuf: "a"},
				{n: 1, wantN: 1, wantBuf: "b"},
				{n: 1, wantN: 1, wantBuf: "c"},
				{n: 1, wantN: 1, wantBuf: "d"},
				{n: 1, wantN: 1, wantBuf: "e"},
				{n: 1, wantN: 0, wantBuf: ""},
			},
		},
		{
			// Verify buf4 leftover is not re-used after EOF.
			// file="ab": read4 returns "ab" (2 chars). read(1) drains "a",
			// leaving "b" in buf4. read(1) drains "b". read(1) triggers
			// a new read4 which returns 0 → EOF.
			name: "buf4_leftover_then_eof",
			file: "ab",
			calls: []call{
				{n: 1, wantN: 1, wantBuf: "a"},
				{n: 1, wantN: 1, wantBuf: "b"},
				{n: 1, wantN: 0, wantBuf: ""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			restore := mockFile(tc.file)
			defer restore()

			r := &Read4{} // fresh instance per test — state resets
			for i, c := range tc.calls {
				buf := make([]byte, c.n+1) // +1 to detect over-writes
				got := r.read(buf, c.n)
				assert.Equal(t, c.wantN, got, "call[%d] n=%d: chars returned", i, c.n)
				assert.Equal(t, c.wantBuf, string(buf[:got]), "call[%d] n=%d: content", i, c.n)
			}
		})
	}
}
