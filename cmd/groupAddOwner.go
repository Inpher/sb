package cmd

import (
	"fmt"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupAddOwner describes the command
type GroupAddOwner struct{}

func init() {
	commands.RegisterCommand("group owner add", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupAddOwner), models.GroupOwner, helpers.Helper{
				Header:      "add an account as a group owner",
				Usage:       "group owner add --account USERNAME --group GROUP",
				Description: "add an account as a group owner",
				Aliases:     []string{"groupAddOwner"},
			}, map[string]commands.Argument{
				"account": {
					Required:    true,
					Description: "The username of the account",
				},
				"group": {
					Required:    true,
					Description: "The group to which attach of the account",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *GroupAddOwner) Checks(ct *commands.Context) error {

	// Check if the user exists
	user, err := models.LoadUser(ct.FormattedArguments["account"])
	if err != nil {
		return fmt.Errorf("account %s doesn't exist", ct.FormattedArguments["account"])
	}
	if user.IsOwnerOfGroup(ct.FormattedArguments["group"]) {
		return fmt.Errorf("account %s is already a group owner", ct.FormattedArguments["account"])
	}

	return nil
}

// Execute executes the command
func (c *GroupAddOwner) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"group":   ct.FormattedArguments["group"],
		"account": ct.FormattedArguments["account"],
	}

	err = c.Replicate(repl)

	return
}

func (c *GroupAddOwner) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *GroupAddOwner) Replicate(repl models.ReplicationData) (err error) {

	err = helpers.AddAccountInGroup(repl["group"], repl["account"], "o")
	if err != nil {
		return
	}

	fmt.Printf("Account %s was successfully added as an owner of group %s\n", repl["account"], repl["group"])

	return
}
