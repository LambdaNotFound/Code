package apidesign

/**
 * The Stateful Buffer Pattern
 *
 * "An external source delivers data in chunks I don't control, but my caller wants a different amount each time."
 *
 * state = {
 *     buf      // leftover from last external call
 *     cursor   // where are we in the external source
 * }
 *
 * Get(n):
 *     1. drain buf first
 *     2. if still need more AND source not exhausted → fetch externally
 *     3. stage fetch into buf (never directly into result)
 *     4. repeat
 */
type ReadBuffer struct {
	buffer        []int
	fetchedOffset int
	readOffset    int

	exhausted bool
	cursor    int
}

type Result struct {
	nextPage *int
	results  []int
}

var fetchPage = func(page int) Result { return Result{} }

func (r *ReadBuffer) fetch(numResult int) []int {
	result := []int{}

	writeCount := 0
	for writeCount < numResult {
		if r.readOffset < r.fetchedOffset { // 1. drain bufffer first
			result = append(result, r.buffer[r.readOffset])
			writeCount++
			r.readOffset++
		} else { // 2. buffer exhausted, refill
			if r.exhausted {
				break
			}

			res := fetchPage(r.cursor)
			if res.nextPage == nil {
				r.exhausted = true
			} else {
				r.cursor = *res.nextPage
			}
			// 3. stage fetch into buf (never directly into result)
			r.buffer = res.results
			r.fetchedOffset = len(res.results)
			r.readOffset = 0
		}
	}

	return result
}

// read4 is the provided API — reads up to 4 chars from the file into buf4
// and returns the number of characters actually read (0-4).
// Declared as a variable so tests can swap in a mock without build-time linkage.
var read4 = func(buf4 []byte) int { return 0 }

/**
 * 158: Read N Characters Given read4 II — Call Multiple Times
 *
 * read4 reads 4 consecutive characters from the file,
 * then writes those characters into the buffer array.
 * The return value is the number of actual characters read.
 *
 */
type Read4 struct {
	buffer        [4]byte
	fetchedOffset int
	readOffset    int
}

func (r *Read4) read(buf []byte, n int) int {
	writeCount := 0

	for writeCount < n {
		if r.readOffset < r.fetchedOffset { // 1. drain bufffer first
			buf[writeCount] = r.buffer[r.readOffset]
			writeCount++
			r.readOffset++
		} else { // 2. buffer exhausted, refill
			// 3. stage fetch into buf (never directly into result)
			r.fetchedOffset = read4(r.buffer[:]) // [:] converts an array to a slice.
			r.readOffset = 0

			if r.fetchedOffset == 0 {
				break // EOF
			}
		}
	}

	return writeCount
}

/**
 * 157: Read N Characters Given Read4
 *
 * no state
 */
func read(buf []byte, n int) int {
	buf4 := make([]byte, 4)
	charsRead := 0

	for charsRead < n {
		count := read4(buf4)
		// copy min(count, n-charsRead) chars
		for i := 0; i < count && charsRead < n; i++ {
			buf[charsRead] = buf4[i]
			charsRead++
		}
		if count < 4 {
			break // EOF
		}
	}
	return charsRead
}
