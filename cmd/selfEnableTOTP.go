package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/inpher/sb/internal/commands"
	"github.com/inpher/sb/internal/config"
	"github.com/inpher/sb/internal/helpers"
	"github.com/inpher/sb/internal/models"

	"github.com/fatih/color"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
)

// SelfEnableTOTP describes the SelfEnableTOTP command
type SelfEnableTOTP struct{}

func init() {
	commands.RegisterCommand("self totp enable", func() (c commands.Command, r models.Right, helper helpers.Helper, args map[string]commands.Argument) {
		return new(SelfEnableTOTP), models.Public, helpers.Helper{
			Header:      "enable TOTP on the account",
			Usage:       "self totp enable",
			Description: "enable TOTP on the account",
			Aliases:     []string{"selfEnableTOTP"},
		}, map[string]commands.Argument{}
	})
}

// Checks checks whether or not the user can execute this method
func (c *SelfEnableTOTP) Checks(ct *commands.Context) error {

	// We're building on top of pam_google_authenticator, let's check the server is setup correctly
	_, err := exec.LookPath("google-authenticator")
	if err != nil {
		return fmt.Errorf("the server is not configured for TOTP")
	}
	return nil
}

// Execute executes the command
func (c *SelfEnableTOTP) Execute(ct *commands.Context) (repl models.ReplicationData, cmdError error, err error) {

	green := color.New(color.FgGreen).SprintFunc()

	enabled, _, _ := ct.User.GetTOTP()
	if enabled {
		fmt.Printf("TOTP is already enabled on this account!")
		return
	}

	// Generate the random emergency codes
	randomCodes := helpers.GetRandomStrings(5, 8)

	// Generate the TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.GetSBHostname(),
		AccountName: fmt.Sprintf("%s@%s", ct.User.User.Username, config.GetSBName()),
	})
	if err != nil {
		return
	}
	// Build a lighter URL than the one provided by the module
	url := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", key.AccountName(), key.Secret(), key.Issuer())

	// Display secret informations
	fmt.Printf("TOTP is about to be enabled on your account.\n")
	fmt.Println()
	fmt.Printf("Here are the information about the generated setup:\n")
	fmt.Printf("  -> otpauth URL: %s\n", green(url))
	fmt.Printf("  -> secret key: %s\n", green(key.Secret()))
	fmt.Printf("  -> QR-Code:\n")
	qrterminal.GenerateWithConfig(
		url,
		qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
			QuietZone: 0,
		},
	)
	fmt.Println()

	// Check that the user is able to provide codes
	validationCode := c.getValidationCode("To validate the setup, please enter a validation code")
	if !totp.Validate(validationCode, key.Secret()) {
		err = fmt.Errorf("the validation code is invalid: TOTP has not been enabled on your account")
		return
	}

	repl = models.ReplicationData{
		"account":      ct.User.User.Username,
		"secret":       key.Secret(),
		"random-codes": strings.Join(randomCodes, ";"),
	}

	err = c.Replicate(repl)
	if err != nil {
		return
	}

	// Display the emergency codes
	fmt.Println()
	fmt.Printf("Here are your %d emergency codes:\n", len(randomCodes))
	for _, str := range randomCodes {
		fmt.Printf("  - %s\n", green(str))
	}
	fmt.Printf("Be sure to store them in a secure place, they will never be displayed again\n")
	fmt.Println()

	// Finally, say it's done!
	fmt.Printf("%s\n", green("TOTP was successfully activated on your account!"))

	return
}

func (c *SelfEnableTOTP) PostExecute(repl models.ReplicationData) (err error) {
	return
}

func (c *SelfEnableTOTP) Replicate(repl models.ReplicationData) (err error) {

	user, err := models.LoadUser(repl["account"])
	if err != nil {
		return
	}

	// Store the secrets (and thus, enable the TOTP on the account)
	err = user.SetTOTPSecret(repl["secret"], strings.Split(repl["random-codes"], ";"))
	if err != nil {
		return fmt.Errorf("an error occurred: TOTP has not been enabled on your account")
	}

	return
}

func (c *SelfEnableTOTP) getValidationCode(prompt string) (validationCode string) {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
