package hexed

import (
	"fmt"
	"io"
)

type Encoder struct {
	w   io.Writer
	buf []byte
	n   int
	e   error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Write(b []byte) (int, error) {
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

func (e *Encoder) Close() error {
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

func (e *Encoder) drain() error {
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

func (e *Encoder) chunks(b []byte) []any {
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
