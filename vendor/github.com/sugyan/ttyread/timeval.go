package ttyread

// TimeVal type
type TimeVal struct {
	Sec  int32
	Usec int32
}

// Add returns sum of TimeVal
func (tv1 TimeVal) Add(tv2 TimeVal) TimeVal {
	sec := tv1.Sec + tv2.Sec
	usec := tv1.Usec + tv2.Usec
	for usec >= 1000000 {
		sec++
		usec -= 1000000
	}
	return TimeVal{
		Sec:  sec,
		Usec: usec,
	}
}

// Subtract returns diff of TimeVal
func (tv1 TimeVal) Subtract(tv2 TimeVal) TimeVal {
	sec := tv1.Sec - tv2.Sec
	usec := tv1.Usec - tv2.Usec
	if usec < 0 {
		sec--
		usec += 1000000
	}
	return TimeVal{
		Sec:  sec,
		Usec: usec,
	}
}
