package cmd

import (
	"fmt"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupDelACLKeeper describes the command
type GroupDelACLKeeper struct{}

func init() {
	commands.RegisterCommand("group acl-keeper remove", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupDelACLKeeper), models.GroupOwner, helpers.Helper{
				Header:      "remove an account from the ACL keepers of a group",
				Usage:       "group acl-keeper delete --account USERNAME --group GROUP",
				Description: "remove an account from the ACL keepers of a group",
				Aliases:     []string{"groupDelACLKeeper"},
			}, map[string]commands.Argument{
				"account": {
					Required:    true,
					Description: "The username of the account",
				},
				"group": {
					Required:    true,
					Description: "The group to which remove the account from",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *GroupDelACLKeeper) Checks(ct *commands.Context) error {

	// Check if the user exists
	user, err := models.LoadUser(ct.FormattedArguments["account"])
	if err != nil {
		return fmt.Errorf("account %s doesn't exist", ct.FormattedArguments["account"])
	}
	if !user.IsACLKeeperOfGroup(ct.FormattedArguments["group"]) {
		return fmt.Errorf("account %s is already not a group ACL keeper", ct.FormattedArguments["account"])
	}

	return nil
}

// Execute executes the command
func (c *GroupDelACLKeeper) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"group":   ct.FormattedArguments["group"],
		"account": ct.FormattedArguments["account"],
	}

	err = c.Replicate(repl)

	return
}

func (c *GroupDelACLKeeper) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *GroupDelACLKeeper) Replicate(repl models.ReplicationData) (err error) {

	err = helpers.RemoveAccountFromGroup(repl["group"], repl["account"], "aclk")
	if err != nil {
		return
	}

	fmt.Printf("Account %s was successfully removed from the ACL keepers of group %s\n", repl["account"], repl["group"])

	return
}
