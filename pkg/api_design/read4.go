package apidesign

// Provided API — reads up to 4 chars from file into buf4
// Returns number of characters actually read (0-4)
// Has its own internal file pointer that advances on each call
func read4(buf4 []byte) int

/**
 * 158: Read N Characters Given read4 II — Call Multiple Times
 *
 * read4 reads 4 consecutive characters from the file,
 * then writes those characters into the buffer array.
 * The return value is the number of actual characters read.
 *
 */
type Read4 struct {
	buf4     [4]byte // internal buffer from read4
	buf4Idx  int     // next unread position in buf4
	buf4Size int     // how many chars read4 returned last fill
}

func (r *Read4) read(buf []byte, n int) int {
	charsRead := 0

	for charsRead < n {
		// refill internal buffer only when fully consumed
		if r.buf4Idx == r.buf4Size {
			r.buf4Size = read4(r.buf4[:]) // [:] converts an array to a slice.
			r.buf4Idx = 0

			if r.buf4Size == 0 {
				break // EOF
			}
		}

		// drain from internal buffer into buf
		buf[charsRead] = r.buf4[r.buf4Idx]
		charsRead++
		r.buf4Idx++
	}

	return charsRead
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
