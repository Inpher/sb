package cmd

import (
	"fmt"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupAddACLKeeper describes the command
type GroupAddACLKeeper struct{}

func init() {
	commands.RegisterCommand("group acl-keeper add", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupAddACLKeeper), models.GroupOwner, helpers.Helper{
				Header:      "add an account as a group ACL keeper",
				Usage:       "group acl-keeper add --account USERNAME --group GROUP",
				Description: "add an account as a group ACL keeper",
				Aliases:     []string{"groupAddACLKeeper"},
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
func (c *GroupAddACLKeeper) Checks(ct *commands.Context) error {

	// Check if the user exists
	user, err := models.LoadUser(ct.FormattedArguments["account"])
	if err != nil {
		return fmt.Errorf("account %s doesn't exist", ct.FormattedArguments["account"])
	}
	if user.IsACLKeeperOfGroup(ct.FormattedArguments["group"]) {
		return fmt.Errorf("account %s is already a group ACL keeper", ct.FormattedArguments["account"])
	}

	return nil
}

// Execute executes the command
func (c *GroupAddACLKeeper) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"group":   ct.FormattedArguments["group"],
		"account": ct.FormattedArguments["account"],
	}

	err = c.Replicate(repl)

	return
}

func (c *GroupAddACLKeeper) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *GroupAddACLKeeper) Replicate(repl models.ReplicationData) (err error) {

	err = helpers.AddAccountInGroup(repl["group"], repl["account"], "aclk")
	if err != nil {
		return
	}

	fmt.Printf("Account %s was successfully added as an ACL keeper of group %s\n", repl["account"], repl["group"])

	return
}
