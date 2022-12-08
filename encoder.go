package hexed

import (
	"fmt"
	"io"
)

type enc struct {
	w   io.Writer
	buf []byte
	n   int
	e   error
}

// NewEncoder returns an encoder object wrapping the given io.Writer. An
// encoder writes hex-editor-style lines to the wrapped io.Writer. For each
// 16-byte chunk of data written to the Encoder, a single line will be written
// to the output. Each line consists of the hexadecimal adress of the first
// byte on the line, followed byte space-separated, hex-encoded byte pairs,
// followed by the string representation of the bytes with non-printable
// characters replaced by '.'.
//
// For example: writing the string 'totally\tradical!' to the Encoder will
// result in the following being written to the underlying io.Writer:
//
// 00000000: 746f 7461 6c6c 7909 7261 6469 6361 6c21  totally.radical!
//
// The caller must Close the encoder to flush any partially written blocks. 
func NewEncoder(w io.Writer) io.WriteCloser {
	return &enc{w: w}
}

// Write implements the io.Writer interface
func (e *enc) Write(b []byte) (int, error) {
	if e.e != nil {
		return 0, e.e
	}
	e.buf = append(e.buf, b...)
	if err := e.drain(); err != nil {
		e.e = err
		return 0, err
	}
	return len(b), nil
}

// Close flushes any remaining data in the buffer. Further writes to the
// encoder will return io.EOF
func (e *enc) Close() error {
	if err := e.drain(); err != nil {
		e.e = err
		return err
	}
	if len(e.buf) == 0 {
		e.e = io.EOF
		return nil
	}

	if _, err := fmt.Fprintf(e.w, fmtStrings[len(e.buf)], e.chunks(e.buf)...); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(e.w, "%s\n", printable(e.buf)); err != nil {
		return err
	}
	return nil
}

func (e *enc) drain() error {
	for len(e.buf) >= 16 {
		if _, err := fmt.Fprintf(e.w, fmtStrings[16], e.chunks(e.buf[:16])...); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(e.w, "%s\n", printable(e.buf[:16])); err != nil {
			e.e = err
			return err
		}
		e.buf = e.buf[16:]
		e.n += 16
	}
	return nil
}

func (e *enc) chunks(b []byte) []any {
	chunks := []any{e.n}
	for i := 0; i < (len(b)/2)*2; i += 2 {
		chunks = append(chunks, b[i:i+2])
	}
	if len(b)%2 == 1 {
		chunks = append(chunks, b[len(b)-1:])
	}
	return chunks
}

func printable(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}

	out := make([]byte, len(b))
	for i, c := range b {
		if c < 32 || c > 127 {
			out[i] = '.'
		} else {
			out[i] = c
		}
	}
	return out
}

var fmtStrings = []string{
	"",
	"%08x: %02x                                       ",
	"%08x: %04x                                     ",
	"%08x: %04x %02x                                  ",
	"%08x: %04x %04x                                ",
	"%08x: %04x %04x %02x                             ",
	"%08x: %04x %04x %04x                           ",
	"%08x: %04x %04x %04x %02x                        ",
	"%08x: %04x %04x %04x %04x                      ",
	"%08x: %04x %04x %04x %04x %02x                   ",
	"%08x: %04x %04x %04x %04x %04x                 ",
	"%08x: %04x %04x %04x %04x %04x %02x              ",
	"%08x: %04x %04x %04x %04x %04x %04x            ",
	"%08x: %04x %04x %04x %04x %04x %04x %02x         ",
	"%08x: %04x %04x %04x %04x %04x %04x %04x       ",
	"%08x: %04x %04x %04x %04x %04x %04x %04x %02x    ",
	"%08x: %04x %04x %04x %04x %04x %04x %04x %04x  ",
}
