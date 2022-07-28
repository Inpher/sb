package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/config"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"
	"github.com/inpher/sb/internal/storage"
	"github.com/pkg/errors"

	"github.com/golgeek/ttyrec2gif"
	"golang.org/x/term"
)

// SelfGetSessionAsGif describes the selfListAccesses command
type SelfGetSessionAsGif struct{}

func init() {
	commands.RegisterCommand("self session gif", func() (c commands.Command, r models.Right, h helpers.Helper, args map[string]commands.Argument) {
		return new(SelfGetSessionAsGif), models.Public, helpers.Helper{
				Header:      "get a recording of an SSH session as a gif",
				Usage:       "self session gif",
				Description: "get a recording of an SSH session as a gif",
				Aliases:     []string{"selfGetSessionAsGif"},
			}, map[string]commands.Argument{
				"session-id": {
					Required:    true,
					Description: "The session recording ID to convert as a GIF",
				},
				"repeat": {
					Required:    false,
					Description: "Specify if animation is repeated",
					Type:        commands.BOOL,
				},
				"speed": {
					Required:     false,
					Description:  "Specify the play speed factor of the session (default is \"1.0\")",
					DefaultValue: "1.0",
				},
			}
	})
}

// Checks checks whether or not the user can execute this method
func (c *SelfGetSessionAsGif) Checks(ct *commands.Context) error {

	_, err := strconv.ParseFloat(ct.FormattedArguments["speed"], 64)
	if err != nil {
		return fmt.Errorf("argument speed is not a valid float")
	}

	return nil
}

// Execute executes the command
func (c *SelfGetSessionAsGif) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	filename := fmt.Sprintf("%s.ttyrec", ct.FormattedArguments["session-id"])
	localFilepath := fmt.Sprintf("%s/%s", ct.User.GetTtyrecDirectory(), filename)
	outputFile := fmt.Sprintf("%s/%s.ttyrec.gif", ct.User.GetTtyrecDirectory(), ct.FormattedArguments["session-id"])

	// If TTYRecs offloading is enabled, we start by getting the ttyrec file from a storage
	ttyRecsOffloadingConfig := config.GetTTYRecsOffloadingConfig()
	if ttyRecsOffloadingConfig.Enabled {

		var rs storage.Storage
		rs, err = storage.GetStorage(ttyRecsOffloadingConfig)
		if err != nil {
			return
		}

		err = rs.GetFromStorage(fmt.Sprintf("%s.bin", filename), fmt.Sprintf("%s.bin", localFilepath))
		if err != nil {
			return
		}

		err = helpers.DecryptFile(fmt.Sprintf("%s.bin", localFilepath), localFilepath, config.GetEncryptionKey())
		if err != nil {
			return
		}

		err = os.Remove(fmt.Sprintf("%s.bin", localFilepath))
		if err != nil {
			return
		}
	}

	_, repeat := ct.FormattedArguments["repeat"]
	speed, err := strconv.ParseFloat(ct.FormattedArguments["speed"], 64)
	if err != nil {
		return
	}

	generator := ttyrec2gif.NewGifGenerator()
	generator.Speed = speed
	generator.NoLoop = !repeat
	err = generator.Generate(localFilepath, outputFile)
	if err != nil {
		err = errors.Wrap(err, "unable to generate GIF from TTYRec")
		return
	}

	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return
	}

	// We set stdout in raw mode to avoid \r\n transformations by ssh -t on client side
	_, err = term.MakeRaw(syscall.Stdout)
	if err != nil {
		return
	}

	fmt.Printf("%s", string(content))

	err = os.Remove(outputFile)
	if err != nil {
		return
	}

	if ttyRecsOffloadingConfig.Enabled {
		err = os.Remove(localFilepath)
		if err != nil {
			return
		}
	}

	return
}

func (c *SelfGetSessionAsGif) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *SelfGetSessionAsGif) Replicate(repl models.ReplicationData) (err error) {
	return
}
