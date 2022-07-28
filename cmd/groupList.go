package cmd

import (
	"fmt"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// GroupList describes the help command
type GroupList struct{}

func init() {
	commands.RegisterCommand("groups list", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(GroupList), models.Public, helpers.Helper{
				Header:      "display the list of groups",
				Usage:       "groups list [--all]",
				Description: "display the list of groups",
				Aliases:     []string{"groupList"},
			}, map[string]commands.Argument{
				"all": {
					Required:    false,
					Description: "List all the groups instead of just the ones you're a member of",
					Type:        commands.BOOL,
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *GroupList) Checks(ct *commands.Context) error {
	// No specific rights needed but a sb account
	return nil
}

// Execute executes the command
func (c *GroupList) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	var groups map[string]*models.Group

	_, ok := ct.FormattedArguments["all"]
	if ok {
		groups, err = models.GetAllSBGroups()
		if err != nil {
			return
		}
	} else {
		groups, err = ct.User.GetAllGroups()
		if err != nil {
			return
		}
	}

	if len(groups) == 0 {
		if ok {
			fmt.Println("No group was created on this sb instance yet")
		} else {
			fmt.Println("You don't have access to any group yet. Use --all to list all sb groups")
		}
	} else {
		fmt.Println("Here are the list of sb groups:")
		for groupName := range groups {
			fmt.Printf("%s\n", groupName)
		}
		if !ok {
			fmt.Println("If you want to see all the groups, even the ones you don't have access to, use --all")
		}
	}

	return
}

func (c *GroupList) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *GroupList) Replicate(repl models.ReplicationData) (err error) {
	return
}
