package ttyrec

import (
	"errors"
	"io"
	"io/ioutil"
)

// ErrReadSeeker is returned if the provided Reader does not provide the io.ReadSeeker interface.
var ErrReadSeeker = errors.New("ttyrec: provided Reader does not implement io.ReadSeeker")

// ErrIllegalSeek is returned if an illegal seek operation is requested, such as seeking to an out of bounds frame.
var ErrIllegalSeek = errors.New("ttyrec: seek to an invalid offset")

// Decoder for TTY recordings.
//
// The decoder methods are not concurrency safe.
type Decoder struct {
	r  io.Reader
	rs io.ReadSeeker

	// started indicates if we have started reading
	started bool

	// startedAt is the first timeVal we see
	startedAt TimeVal

	// sequence of the current frame
	sequence int

	// offset for the decoded frames
	offset int64
	chunks []int64
}

// NewDecoder returns a new Decoder for the provided Reader.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r: r,
	}
	if rs, ok := r.(io.ReadSeeker); ok {
		d.rs = rs
	}
	return d
}

// DecodeFrame decodes a single frame.
func (d *Decoder) DecodeFrame() (*Frame, error) {
	return d.decodeFrame(false)
}

func (d *Decoder) decodeFrame(discard bool) (*Frame, error) {
	var (
		f   Frame
		n   int64
		err error
	)

	// Read header.
	if n, err = f.Header.ReadFrom(d.r); err != nil {
		return nil, err
	}

	// Read data.
	if f.Len > 0 {
		n += int64(f.Len)
		l := io.LimitReader(d.r, int64(f.Len))
		if discard {
			if _, err = io.Copy(ioutil.Discard, l); err != nil {
				return nil, err
			}
		} else if f.Data, err = ioutil.ReadAll(l); err != nil {
			return nil, err
		}
	}

	// Record first time stamp, the rest of the frame times are relative to the first.
	if !d.started {
		d.started = true
		d.startedAt = f.Header.Time
	}

	// Bookkeeping, tracking the sequence number and size of the chunks.
	d.sequence++
	if len(d.chunks) < d.sequence {
		d.chunks = append(d.chunks, n)
	}

	return &f, nil
}

// StopFunc is used to interrupt a stream.
type StopFunc func()

// DecodeStream returns a stream of frames.
func (d *Decoder) DecodeStream() (<-chan *Frame, StopFunc) {
	var (
		output = make(chan *Frame)
		stop   = make(chan struct{})
	)

	go func(output chan<- *Frame, stop chan struct{}) {
		defer close(output)

		for {
			f, err := d.decodeFrame(false)
			if err != nil {
				return
			}
			select {
			case <-stop:
				return
			case output <- f:
			}
		}
	}(output, stop)

	return output, func() {
		close(stop)
	}
}

// Frame returns the current frame number.
func (d *Decoder) Frame() int {
	return d.sequence
}

// SeekToFrame seeks to the specified frame offset. Whence can be any of the
// io.SeekStart (relative to first frame) or io.SeekCurrent (relative to
// current frame). io.SeekEnd is not supported.
func (d *Decoder) SeekToFrame(offset, whence int) error {
	var n int
	switch whence {
	case io.SeekStart:
		n = offset
	case io.SeekCurrent:
		n = d.sequence + offset
	default:
		return ErrIllegalSeek
	}

	if n < 0 {
		return ErrIllegalSeek
	}
	if delta := n - d.sequence; delta < 0 {
		return d.rewindFrames(-delta)
	} else if delta > 0 {
		return d.advanceFrames(delta)
	}

	return nil
}

func (d *Decoder) advanceFrames(n int) error {
	for i := 0; i < n; i++ {
		if _, err := d.decodeFrame(true); err != nil {
			return err
		}
	}
	return nil
}

/*
func (d *Decoder) rewindFrame() error {
	if d.rs == nil {
		return ErrReadSeeker
	} else if d.sequence == 0 {
		return ErrIllegalSeek
	}
	if _, err := d.rs.Seek(-d.chunks[d.sequence-1], io.SeekCurrent); err != nil {
		return err
	}
	d.sequence--
	return nil
}

func (d *Decoder) rewindFrames(n int) error {
	for ; n > 0; n-- {
		if err := d.rewindFrame(); err != nil {
			return err
		}
	}
	return nil
}
*/

func (d *Decoder) rewindFrames(n int) error {
	var offset int64
	for i := 0; i < n; i++ {
		offset -= d.chunks[d.sequence-1]
		d.sequence--
	}
	if _, err := d.rs.Seek(offset, io.SeekCurrent); err != nil {
		return err
	}
	return nil
}
