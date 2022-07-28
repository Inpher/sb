package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Helper describes the basic properties of a sb Helper type
type Helper struct {
	Header      string
	Usage       string
	Description string
	Aliases     []string
}

// ParseArguments parses the os.Args arguments
func ParseArgumentsNew(clArgs []string) (c string, ca []string, rest []string, err error) {

	// We initialize our client with "ssh"
	c = "ssh"

	// Absolutely no arguments were provided (that's odd, we should at least have ourselves at $0)
	if len(clArgs) == 0 {
		return
	}

	// We start by dropping the script name
	clArgs = clArgs[1:]

	// Absolutely no arguments were provided
	if len(clArgs) == 0 {
		return
	}

	// We might have been called from two ways:
	// - directly, then arguments are: ["command", "--arg1", "val1, "--arg2", ...]
	// - through SSH client, then arguments are: ["-c", "command --arg1 val1 --arg2 ..."]
	if clArgs[0] != "-c" {
		// Nothing to do
	} else {
		if len(clArgs) == 1 {
			// No arguments were provided
			return
		}
		// We drop the "-c" and convert the arguments string to an array of arguments
		clArgs, err = ParseCommandLine(clArgs[1])
		if err != nil {
			return
		}
	}

	// We now have two different cases:
	// - just a basic slice of arguments (in case of direct call and SSH)
	// - something starting with ["mosh-server", "new", ..., "--", args...]
	if clArgs[0] == "mosh-server" {
		c = "mosh"

		// We drop the non-arguments part of the mosh command line.
		// Second argument should be "new", in case it is not, we check this case
		if clArgs[1] == "new" {
			clArgs = clArgs[2:]
		} else {
			clArgs = clArgs[1:]
		}

		// Mosh wraps around every argument with ' and escapes the ', we'll try and fix that
		for i := 0; i < len(clArgs); i++ {
			clArgs[i] = strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(clArgs[i], `\'`, "'"), "'"), "'")
		}

		// Let's build a flag parser for mosh options!
		moshFlagSet := flag.NewFlagSet("mosh", flag.ContinueOnError)

		// The flag package displays a nice message if an undeclared flag is found... but that's not what we want!
		// Let's redirect output to an abandoned buffer
		var buf bytes.Buffer
		moshFlagSet.SetOutput(&buf)

		// This flags are extracted from mosh-server man page
		iface := moshFlagSet.Bool("s", false, "bind to the local interface used for an incoming SSH connection, given in the SSH_CONNECTION environment variable (for multihomed hosts)")
		v := moshFlagSet.Bool("v", false, "Print some debugging information even after detaching.  More instances of this flag will result in more debugging information.")
		ip := moshFlagSet.String("i", "", "IP address of the local interface to bind (for multihomed hosts)")
		colors := moshFlagSet.String("c", "8", "Number of colors to advertise to applications through TERM (e.g. 8, 256)")
		lang := moshFlagSet.String("l", "", "Locale-related environment variable to try as part of a fallback environment, if the startup environment does not specify a character set of UTF-8.")
		moshFlagSet.String("p", "", "UDP port number or port-range to bind.  -p 0 will let the operating system pick an available UDP port.")

		// Let's parse our mosh-server arguments
		moshFlagSet.Parse(clArgs)

		// And keep everything that was trailing for the next step
		clArgs = moshFlagSet.Args()

		// We push everything we parsed into our clientArguments slice
		// Everything except the ports to use, as they're coded in our configuration file
		ca = make([]string, 0)
		if *iface {
			ca = append(ca, "-s")
		}
		if *v {
			ca = append(ca, "-v")
		}
		if *ip != "" {
			ca = append(ca, "-i", *ip)
		}
		if *colors != "" {
			ca = append(ca, "-c", *colors)
		}
		if *lang != "" {
			ca = append(ca, "-l", *lang)
		}
	}

	// We return all remaining arguments, that's what the user really wanted to give us
	rest = clArgs

	return
}

// ParseArguments parses the os.Args arguments
func ParseArguments(clArgs []string) (c string, ca []string, ba map[string]bool, rest []string, err error) {

	// We initialize our client with "ssh"
	c = "ssh"

	// Absolutely no arguments were provided (that's odd, we should at least have ourselves at $0)
	if len(clArgs) == 0 {
		return
	}

	// We start by dropping the script name
	clArgs = clArgs[1:]

	// Absolutely no arguments were provided
	if len(clArgs) == 0 {
		return
	}

	// We might have been called from two ways:
	// - directly, then arguments are: ["command", "--arg1", "val1, "--arg2", ...]
	// - through SSH client, then arguments are: ["-c", "command --arg1 val1 --arg2 ..."]
	if clArgs[0] != "-c" {
		// Nothing to do
	} else {
		if len(clArgs) == 1 {
			// No arguments were provided
			return
		}
		// We drop the "-c" and convert the arguments string to an array of arguments
		clArgs, err = ParseCommandLine(clArgs[1])
		if err != nil {
			return
		}
	}

	// We now have two different cases:
	// - just a basic slice of arguments (in case of direct call and SSH)
	// - something starting with ["mosh-server", "new", ..., "--", args...]
	if clArgs[0] == "mosh-server" {
		c = "mosh"

		// We drop the non-arguments part of the mosh command line.
		// Second argument should be "new", in case it is not, we check this case
		if clArgs[1] == "new" {
			clArgs = clArgs[2:]
		} else {
			clArgs = clArgs[1:]
		}

		// Mosh wraps around every argument with ' and escapes the ', we'll try and fix that
		for i := 0; i < len(clArgs); i++ {
			clArgs[i] = strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(clArgs[i], `\'`, "'"), "'"), "'")
		}

		// Let's build a flag parser for mosh options!
		moshFlagSet := flag.NewFlagSet("mosh", flag.ContinueOnError)

		// The flag package displays a nice message if an undeclared flag is found... but that's not what we want!
		// Let's redirect output to an abandoned buffer
		var buf bytes.Buffer
		moshFlagSet.SetOutput(&buf)

		// This flags are extracted from mosh-server man page
		iface := moshFlagSet.Bool("s", false, "bind to the local interface used for an incoming SSH connection, given in the SSH_CONNECTION environment variable (for multihomed hosts)")
		v := moshFlagSet.Bool("v", false, "Print some debugging information even after detaching.  More instances of this flag will result in more debugging information.")
		ip := moshFlagSet.String("i", "", "IP address of the local interface to bind (for multihomed hosts)")
		colors := moshFlagSet.String("c", "8", "Number of colors to advertise to applications through TERM (e.g. 8, 256)")
		lang := moshFlagSet.String("l", "", "Locale-related environment variable to try as part of a fallback environment, if the startup environment does not specify a character set of UTF-8.")
		moshFlagSet.String("p", "", "UDP port number or port-range to bind.  -p 0 will let the operating system pick an available UDP port.")

		// Let's parse our mosh-server arguments
		moshFlagSet.Parse(clArgs)

		// And keep everything that was trailing for the next step
		clArgs = moshFlagSet.Args()

		// We push everything we parsed into our clientArguments slice
		// Everything except the ports to use, as they're coded in our configuration file
		ca = make([]string, 0)
		if *iface {
			ca = append(ca, "-s")
		}
		if *v {
			ca = append(ca, "-v")
		}
		if *ip != "" {
			ca = append(ca, "-i", *ip)
		}
		if *colors != "" {
			ca = append(ca, "-c", *colors)
		}
		if *lang != "" {
			ca = append(ca, "-l", *lang)
		}
	}

	clArgs = RegroupCommandArguments(clArgs)

	// Here, we will introduce the parsing of sb arguments
	sbFlagSet := flag.NewFlagSet("sb", flag.ContinueOnError)

	// The flag package displays a nice message if an undeclared flag is found... but that's not what we want!
	// Let's redirect output to an abandoned buffer
	var buf bytes.Buffer
	sbFlagSet.SetOutput(&buf)

	// Let's add our sb flags
	v := sbFlagSet.Bool("v", false, "Debug")
	i := sbFlagSet.Bool("i", false, "Interactive mode")
	d := sbFlagSet.Bool("d", false, "Daemon")

	// Let's parse our sb arguments
	sbFlagSet.Parse(clArgs)

	// And keep everything that was trailing for the next step
	clArgs = sbFlagSet.Args()

	if *v || *i || *d {
		ba = make(map[string]bool)
	}
	if *v {
		ba["verbose"] = true
	}
	if *i {
		ba["interactive"] = true
	}
	if *d {
		ba["daemon"] = true
	}

	// We return all remaining arguments, that's what the user really wanted to give us
	rest = clArgs

	return
}

func RegroupCommandArguments(clArgs []string) (args []string) {

	var j int
	firstArg := clArgs[0]
	for j = 1; j < len(clArgs); j++ {

		arg := clArgs[j]

		if strings.HasPrefix(arg, "-") {
			break
		}

		firstArg = fmt.Sprintf("%s %s", firstArg, arg)
	}
	args = append([]string{firstArg}, clArgs[j:]...)

	return
}

// ParseCommandLine parses the string passed to us by SSH to an array of args
func ParseCommandLine(cmd string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(cmd); i++ {
		c := cmd[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", cmd)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

// GetRandomStrings returns x random strings of y characters
func GetRandomStrings(quantity int, length int) (rdm []string) {

	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("0123456789")

	for i := 0; i < quantity; i++ {

		b := make([]rune, length)
		for j := range b {
			b[j] = letterRunes[rand.Intn(len(letterRunes))]
		}
		rdm = append(rdm, string(b))
	}

	return
}

func GetHostname() (hostname string, err error) {
	return os.Hostname()
}

func DecryptFile(filepathIn, filepathOut, decryptionKey string) (err error) {

	var file *os.File
	var outfile *os.File

	file, err = os.Open(filepathIn)
	if err != nil {
		return
	}
	defer file.Close()

	block, err := aes.NewCipher([]byte(decryptionKey))
	if err != nil {
		return
	}

	fi, err := file.Stat()
	if err != nil {
		return
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = file.ReadAt(iv, msgLen)
	if err != nil {
		return
	}

	outfile, err = os.OpenFile(filepathOut, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}
	}

	return
}

func EncryptFile(filepathIn, filepathOut, replicationEncryptionKey string) (err error) {

	var file *os.File
	var outfile *os.File

	file, err = os.Open(filepathIn)
	if err != nil {
		return
	}
	defer file.Close()

	cipherKey := replicationEncryptionKey

	block, err := aes.NewCipher([]byte(cipherKey))
	if err != nil {
		return
	}

	iv := make([]byte, block.BlockSize())
	if _, err = io.ReadFull(cryptorand.Reader, iv); err != nil {
		return
	}

	outfile, err = os.OpenFile(filepathOut, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer outfile.Close()

	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}
	}
	outfile.Write(iv)

	return
}
