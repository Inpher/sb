package ttyrec

import (
	"encoding/binary"
	"io"
	"time"
)

const (
	timeValLen = 8
	headerLen  = timeValLen + 4
)

// ByteOrder is our encoding byte order.
var byteOrder = binary.LittleEndian

// Frame for a recording.
type Frame struct {
	Header
	Data []byte
}

// Header for a Frame.
type Header struct {
	Time TimeVal
	Len  uint32
}

// ReadFrom reads the header values from the provided Reader.
func (h *Header) ReadFrom(r io.Reader) (int64, error) {
	return headerLen, binary.Read(r, byteOrder, h)
}

// WriteTo writes the header values to the provided Writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	return headerLen, binary.Write(w, byteOrder, h)
}

// TimeVal is a struct timeval.
type TimeVal struct {
	Seconds      int32
	MicroSeconds int32
}

// Sub subtracts x from t, returning the difference.
func (t TimeVal) Sub(x TimeVal) time.Duration {
	var (
		ds = time.Duration(t.Seconds) - time.Duration(x.Seconds)
		dµ = time.Duration(t.MicroSeconds) - time.Duration(x.MicroSeconds)
	)
	return ds*time.Second + dµ*time.Microsecond
}

// Set the values based on the provided duration. Negative durations are ignored.
func (t *TimeVal) Set(d time.Duration) {
	if d < 0 {
		return
	}
	// nano -> micro
	d /= time.Microsecond
	t.Seconds = int32(d / 1000000)
	t.MicroSeconds = int32(d - 1000000*time.Duration(t.Seconds))
}
