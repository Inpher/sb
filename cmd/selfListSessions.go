package cmd

import (
	"fmt"
	"strings"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// SelfListSessions describes the selfListAccesses command
type SelfListSessions struct{}

func init() {
	commands.RegisterCommand("self sessions list", func() (c commands.Command, r models.Right, h helpers.Helper, args map[string]commands.Argument) {
		return new(SelfListSessions), models.Public, helpers.Helper{
			Header:      "list your last 20 SSH sessions",
			Usage:       "self sessions list",
			Description: "list your last 20 SSH sessions",
			Aliases:     []string{"selfListSessions"},
		}, map[string]commands.Argument{}
	})
}

// Checks checks whether or not the user can execute this method
func (c *SelfListSessions) Checks(ct *commands.Context) error {
	// No specific rights needed but a sb account
	return nil
}

// Execute executes the command
func (c *SelfListSessions) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	lastSessions, err := ct.User.GetLastSSHSessions(20)
	if err != nil {
		return
	}

	sessions := make([]string, 0, len(lastSessions))
	for id, session := range lastSessions {
		sessions = append(sessions, fmt.Sprintf("%02d: %s", id+1, session.String()))
	}

	if len(sessions) > 0 {
		fmt.Printf(
			"Here is the list of your last 20 SSH sessions:\n%s\n",
			strings.Join(sessions, "\n"),
		)
	} else {
		fmt.Println("You currently don't have any recorded SSH session")
	}

	return
}

func (c *SelfListSessions) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *SelfListSessions) Replicate(repl models.ReplicationData) (err error) {
	return
}
