package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupAddAccess describes the help command
type GroupAddAccess struct{}

func init() {
	commands.RegisterCommand("group access add", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupAddAccess), models.GroupACLKeeper, helpers.Helper{
				Header:      "add a group access to a distant host",
				Usage:       "group access add --group GROUP-NAME --host HOST --user USER [--port PORT --alias ALIAS]",
				Description: "add a group access to a distant host",
				Aliases:     []string{"groupAddAccess"},
			}, map[string]commands.Argument{
				"group": {
					Required:    true,
					Description: "The group name you want to add an access for",
				},
				"host": {
					Required:    true,
					Description: "An IP, IP range or hostname you're granting access to",
				},
				"user": {
					Required:    true,
					Description: "The user you're granting access to",
				},
				"port": {
					Required:     false,
					Description:  "An optional port you're granting access to. If not provided, the default SSH port (22) will be used.",
					DefaultValue: "22",
				},
				"alias": {
					Required:    false,
					Description: "An optional alias to this access (to enable quick access by typing 'sb alias' or 'sb user@alias')",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *GroupAddAccess) Checks(ct *commands.Context) error {

	for commandName := range commands.GetCommands() {
		if strings.EqualFold(ct.FormattedArguments["alias"], commandName) ||
			strings.EqualFold(ct.FormattedArguments["host"], commandName) {
			return fmt.Errorf("the host or alias provided matches a sb command name, please provide a different value")
		}
	}

	return nil
}

func (c *GroupAddAccess) PostExecute(repl models.ReplicationData) (err error) {
	return
}

// Execute executes the command
func (c *GroupAddAccess) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"group":   ct.Group.Name,
		"host":    ct.FormattedArguments["host"],
		"user":    ct.FormattedArguments["user"],
		"port":    ct.FormattedArguments["port"],
		"alias":   ct.FormattedArguments["alias"],
		"comment": fmt.Sprintf("Access granted by %s on %s", ct.User.User.Username, time.Now().Format(time.RFC3339)),
	}

	err = c.Replicate(repl)

	return
}

func (c *GroupAddAccess) Replicate(repl models.ReplicationData) (err error) {

	grp, err := models.GetGroup(repl["group"])
	if err != nil {
		return
	}

	ba, err := grp.AddAccess(
		repl["host"],
		repl["user"],
		repl["port"],
		repl["alias"],
		repl["comment"],
	)
	if err != nil {
		return
	}

	fmt.Printf("The access has been granted: %s\n", ba.String())

	return
}
