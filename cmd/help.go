package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/config"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
)

// Help describes the help command
type Help struct {
}

func init() {
	commands.RegisterCommand("help", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(Help), models.Public, helpers.Helper{
			Header:      "display this help",
			Usage:       "help",
			Description: "display this help",
		}, map[string]commands.Argument{}
	})
}

// Checks checks whether or not the user can execute this method
func (c *Help) Checks(ct *commands.Context) error {
	// No specific rights needed but a sb account
	return nil
}

// Execute executes the command
func (c *Help) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	cmds := commands.GetCommands()

	maxLength := 0
	commandNames := make([]string, 0, len(cmds))

	for commandName := range cmds {

		if !commands.IsAPublicCommand(commandName) {
			continue
		}

		if len(commandName) > maxLength {
			maxLength = len(commandName)
		}
		commandNames = append(commandNames, commandName)
	}
	sort.Strings(commandNames)

	fmt.Printf("Usage: %s [OPTION | HOST | COMMAND]\n", config.GetSBName())
	fmt.Println()
	fmt.Printf("Available options:\n")
	fmt.Printf("  -i: launch sb in interactive mode\n")
	fmt.Println()
	fmt.Printf("Host supported formats:\n")
	fmt.Printf("  - full formats:\n")
	fmt.Printf("    - user@example.com:22\n")
	fmt.Printf("    - user@127.0.0.1:22\n")
	fmt.Printf("  - short formats*:\n")
	fmt.Printf("    - user@example.com : port will be retrieved from granted access\n")
	fmt.Printf("    - example.com:22   : user will be retrieved from granted access\n")
	fmt.Printf("    - example.com      : port and user will be retrieved from granted access\n")
	fmt.Printf("  - alias*:\n")
	fmt.Printf("    - user@alias:port  : host will be retrieved from granted access\n")
	fmt.Printf("    - user@alias       : host and port will be retrieved from granted access\n")
	fmt.Printf("    - alias:port       : host and user will be retrieved from granted access\n")
	fmt.Printf("    - alias            : host, user and port will be retrieved from granted access\n")
	fmt.Printf("* If multiple granted access match a short format or an alias,\n")
	fmt.Printf("user will be interactively prompted to choose the desired access\n")
	fmt.Println()
	fmt.Println("Available commands:")
	for _, commandName := range commandNames {
		_, rightLevel, helper, _ := cmds[commandName]()

		if rightLevel >= models.Private {
			continue
		}

		if rightLevel < models.SBOwner || ct.User.IsOwnerOfGroup("owners") {
			fmt.Printf("  - %-"+strconv.Itoa(maxLength)+"s : %s\n", commandName, helper.Header)
		}
	}

	return
}

func (c *Help) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *Help) Replicate(repl models.ReplicationData) (err error) {
	return
}
