package cmd

import (
	"fmt"
	"time"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// DelAccount describes the command
type DelAccount struct{}

func init() {
	commands.RegisterCommand("account delete", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(DelAccount), models.SBOwner, helpers.Helper{
				Header:      "delete an account from sb",
				Usage:       "account delete --account USERNAME",
				Description: "delete an account from sb",
				Aliases:     []string{"delAccount"},
			}, map[string]commands.Argument{
				"account": {
					Required:    true,
					Description: "The username of the account",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *DelAccount) Checks(ct *commands.Context) error {

	// Check the user is actually a genuine sb user

	users, err := models.GetAllSBUsers()
	if err != nil {
		return err
	}

	found := false
	for _, user := range users {
		if user == ct.FormattedArguments["account"] {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("account doesn't exist")
	}

	return nil
}

// Execute executes the command
func (c *DelAccount) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"account":        ct.FormattedArguments["account"],
		"archive-suffix": fmt.Sprintf("bak_%d", time.Now().Unix()),
	}

	err = c.Replicate(repl)

	return
}

func (c *DelAccount) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *DelAccount) Replicate(repl models.ReplicationData) (err error) {

	err = helpers.DeleteAccount(repl["account"], repl["archive-suffix"])
	if err != nil {
		return
	}

	fmt.Printf("Account %s was successfully deleted\n", repl["account"])

	return
}
