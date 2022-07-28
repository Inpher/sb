package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// SelfAddIngressKey describes the selfAddIngresKey command
type SelfAddIngressKey struct{}

func init() {
	commands.RegisterCommand("self ingress-key add", func() (c commands.Command, r models.Right, h helpers.Helper, args map[string]commands.Argument) {
		return new(SelfAddIngressKey), models.Public, helpers.Helper{
				Header:      "add a new public ingress key (you -> sb) to your account",
				Usage:       "self ingress-key add [--public-key 'KEY']",
				Description: "add a new public ingress key (you -> sb) to your account",
				Aliases:     []string{"selfAddIngressKey"},
			}, map[string]commands.Argument{
				"public-key": {
					Required:    false,
					Description: "Your new public SSH key to deposit on sb (you will need to '\"double escape it\"'); if not present, you'll be prompted for it",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *SelfAddIngressKey) Checks(ct *commands.Context) error {
	// No specific rights needed but a sb account
	return nil
}

// Execute executes the command
func (c *SelfAddIngressKey) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	str, keys, _ := ct.User.DisplayPubKeys("ingress")

	var pk *helpers.PublicKey

	if ct.FormattedArguments["public-key"] == "" {
		// Interactive mode

		// We start by displaying the current keys (output of selfListIngressKeys)
		fmt.Printf("Here is the list of your current ingress public SSH keys (you -> sb):\n%s\n", str)

		// Then we request the new key to add
		for pk == nil {
			fmt.Print("Please paste the public key you want to add: ")
			// Scan stdin to get the key to delete
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			ct.FormattedArguments["public-key"] = scanner.Text()

			pk, err = helpers.CheckStringPK(ct.FormattedArguments["public-key"], keys)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}
	} else {
		// Non-interactive mode

		// We try to validate the provided public-key
		pk, err = helpers.CheckStringPK(ct.FormattedArguments["public-key"], keys)
		if err != nil {
			return
		}
	}

	repl = models.ReplicationData{
		"account":    ct.User.User.Username,
		"public-key": pk.String(),
	}

	err = c.Replicate(repl)

	// We finally display the account's keys (output of selfListIngressKeys)
	str, _, _ = ct.User.DisplayPubKeys("ingress")
	fmt.Printf("Here is the list of your current ingress public SSH keys (you -> sb):\n%s\n", str)

	return
}

func (c *SelfAddIngressKey) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *SelfAddIngressKey) Replicate(repl models.ReplicationData) (err error) {

	user, err := models.LoadUser(repl["account"])
	if err != nil {
		return
	}

	// We add the key to the account
	err = user.AddIngressKey(repl["public-key"])
	if err != nil {
		return
	}

	fmt.Println("Your key was successfully added.")

	return
}
