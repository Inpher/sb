package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// SelfForgetHostkey describes the SelfForgetHostkey command
type SelfForgetHostkey struct{}

func init() {
	commands.RegisterCommand("self hostkey forget", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(SelfForgetHostkey), models.Public, helpers.Helper{
				Header:      "forget a hostkey",
				Usage:       "self hostkey forget",
				Description: "forget a hostkey from your known_hosts file",
			}, map[string]commands.Argument{
				"hostkey": {
					Required:    true,
					Description: "The hostkey to forget",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *SelfForgetHostkey) Checks(ct *commands.Context) error {
	return nil
}

// Execute executes the command
func (c *SelfForgetHostkey) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	repl = models.ReplicationData{
		"account": ct.User.User.Username,
		"hostkey": ct.FormattedArguments["hostkey"],
	}

	err = c.Replicate(repl)
	if err != nil {
		return
	}

	// Finally, say it's done!
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s\n", green("The provided hostkey was successfully forgotten!"))

	return
}

func (c *SelfForgetHostkey) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *SelfForgetHostkey) Replicate(repl models.ReplicationData) (err error) {

	user, err := models.LoadUser(repl["account"])
	if err != nil {
		return
	}

	// Store the secrets (and thus, enable the TOTP on the account)
	err = helpers.RemoveHostKey(repl["account"], user.GetKnownHostsFilepath(), repl["hostkey"])
	if err != nil {
		return fmt.Errorf("an error occurred: unable to forget hostkey: %s", err)
	}

	return
}
