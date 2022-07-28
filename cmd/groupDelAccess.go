package cmd

import (
	"fmt"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupDelAccess describes the help command
type GroupDelAccess struct{}

func init() {
	commands.RegisterCommand("group access remove", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupDelAccess), models.GroupACLKeeper, helpers.Helper{
				Header:      "remove a group access to a distant host",
				Usage:       "group access delete --group GROUP-NAME --host HOST --user USER --port PORT",
				Description: "remove a group access to a distant host",
				Aliases:     []string{"groupDelAccess"},
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
					Required:    true,
					Description: "The port you're granting access to",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *GroupDelAccess) Checks(ct *commands.Context) error {
	return nil
}

// Execute executes the command
func (c *GroupDelAccess) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"group": ct.Group.Name,
		"host":  ct.FormattedArguments["host"],
		"user":  ct.FormattedArguments["user"],
		"port":  ct.FormattedArguments["port"],
	}

	err = c.Replicate(repl)

	return
}

func (c *GroupDelAccess) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *GroupDelAccess) Replicate(repl models.ReplicationData) (err error) {

	group, err := models.GetGroup(repl["group"])
	if err != nil {
		return
	}

	ba, err := group.DeleteAccess(
		repl["host"],
		repl["user"],
		repl["port"],
	)
	if err != nil {
		return
	}

	fmt.Printf("The access %s has been removed\n", ba.ShortString())

	return
}
