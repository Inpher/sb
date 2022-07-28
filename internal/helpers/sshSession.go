package helpers

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// SSHSession represents an SSH session
type SSHSession struct {
	UniqID    string
	StartDate time.Time
	EndDate   time.Time
	UserFrom  string
	IPFrom    string
	PortFrom  string
	HostFrom  string
	HostTo    string
	PortTo    string
	UserTo    string
	Allowed   bool
}

func (s *SSHSession) String() (str string) {

	allowed := "denied"
	var duration string
	sessionEnd := "????-??-?? ??:??:??"
	if !s.EndDate.IsZero() {
		d := s.EndDate.Sub(s.StartDate)
		duration = time.Time{}.Add(d).Format("15:04:05")
		sessionEnd = s.EndDate.Format("2006-01-02 15:04:05")
	} else {
		d := time.Since(s.StartDate)
		duration = fmt.Sprintf("%s (running)", time.Time{}.Add(d).Format("15:04:05"))
	}

	if s.Allowed {
		allowed = "allowed"
	}

	str = fmt.Sprintf(
		`Session ID: %s (%s | %s - %s)
	- From: %s@%s:%s
	- To: %s@%s:%s
	- Duration: %s`,
		s.UniqID, allowed, s.StartDate.Format("2006-01-02 15:04:05"), sessionEnd,
		s.UserFrom, s.HostFrom, s.PortFrom,
		s.UserTo, s.HostTo, s.PortTo,
		duration,
	)

	if s.Allowed {
		return color.New(color.FgGreen).SprintFunc()(str)
	}

	return color.New(color.FgRed).SprintFunc()(str)
}
