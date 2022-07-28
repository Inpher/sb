package ttyrec

import (
	"io"
	"time"
)

// Encoder can write chunks of bytes in a ttyrec format.
type Encoder struct {
	w io.Writer

	// started indicates if we have started writing
	started bool

	// startedAt is the time of first write
	startedAt time.Time
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (e *Encoder) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	header := Header{Len: uint32(len(p))}
	if !e.started {
		e.started = true
		e.startedAt = time.Now()
	} else {
		header.Time.Set(time.Since(e.startedAt))
	}

	// Write header.
	if _, err := header.WriteTo(e.w); err != nil {
		return 0, err
	}

	// Write data.
	return e.w.Write(p)
}
