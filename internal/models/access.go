package models

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/inpher/sb/internal/helpers"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

// Access descibes the basic properties of this struct
type Access struct {
	UniqID  string `gorm:"PRIMARY_KEY"`
	Host    string `gorm:"type:varchar(100);unique_index:host_user_prefix_port"`
	Prefix  string `gorm:"type:varchar(50);unique_index:host_user_prefix_port"`
	Alias   string `gorm:"type:varchar(100);unique_index:host_user_prefix_port"`
	User    string `gorm:"type:varchar(50);unique_index:host_user_prefix_port"`
	Port    int    `gorm:"type:varchar(5);unique_index:host_user_prefix_port"`
	Comment string `gorm:"type:text"`
	IP      net.IP `gorm:"-"`
}

// BeforeCreate will set a UUID if not present
func (ba *Access) BeforeCreate(tx *gorm.DB) (err error) {
	if ba.UniqID == "" {
		ba.UniqID = uuid.New().String()
	}
	return nil
}

// AccessesByKeys describes the basic properties of the struct
type AccessesByKeys struct {
	Keys     []*helpers.SSHKeyPair
	Accesses []*Access
	Type     string
	Group    string
}

// Source describes the basic properties of the struct
type Source struct {
	Type  string
	Group string
}

// Info describes the basic properties of the struct
type Info struct {
	KeyFilepathes []string
	Sources       []*Source
	Authorized    bool
	Accesses      []*Access
}

// BuildSBAccess builds a new sb access to be stored
// It will resolve DNS on host and store prefix
func BuildSBAccess(host, user, port, alias string, strictHostCheck bool) (ba *Access, err error) {

	if host == "" {
		return ba, fmt.Errorf("host is missing")
	}

	var intPort int
	if port != "" {
		intPort, err = strconv.Atoi(port)
		if err != nil {
			return ba, fmt.Errorf("port is invalid")
		}
	}

	typeGiven := ""

	// Is it an IP?
	hostIP := net.ParseIP(host)
	if hostIP != nil {
		// If it is, we'll add the / to convert it into a net.IPNet
		if isIPv4(host) {
			host = fmt.Sprintf("%s/32", host)
		} else {
			host = fmt.Sprintf("%s/128", host)
		}
		typeGiven = "IP"
	}

	// Is it a network? (if it was an IP, it should now be a network)
	hostIP, hostIPNet, errCIDR := net.ParseCIDR(host)
	if errCIDR == nil {
		if typeGiven != "IP" {
			typeGiven = "IPNET"
		}
	} else {

		typeGiven = "HOST"

		// OK, maybe it's a host, we will try to resolve it
		ips, err := net.LookupIP(host)
		switch {
		case err == nil:
			for _, ip := range ips {

				// Convert net.IP as net.IPNet
				ipRange := ip.String()
				if isIPv4(ipRange) {
					ipRange = fmt.Sprintf("%s/32", ipRange)
				} else {
					ipRange = fmt.Sprintf("%s/128", ipRange)
				}

				hostIP, hostIPNet, err = net.ParseCIDR(ipRange)
				if err != nil {
					return ba, err
				}

			}
		case strictHostCheck:
			return ba, fmt.Errorf("host is neither an IP, a prefix or a resolvable host")
		default:
			// We're not strictly checking that the host is valid,
			// so let's consider this is an alias
		}
	}

	var hostIPNetstr string
	if hostIPNet != nil {
		hostIPNetstr = hostIPNet.String()
	}

	if typeGiven == "IP" || typeGiven == "IPNET" {

		// Compute the required /slash to do a reverse lookup
		slashRange := "/32"
		if !isIPv4(hostIPNetstr) {
			slashRange = "/128"
		}

		// We can try and do a reverse lookup to get the host name from the IP (or the IPNet which might be a ipv4/32 or ipv6/128)
		if strings.HasSuffix(hostIPNetstr, slashRange) {

			// If we can't find it, we will store the IP as host
			names, errRecord := net.LookupAddr(hostIP.String())
			if errRecord != nil || len(names) == 0 {
				// We won't throw an error here, it's just a nice to have
				host = hostIP.String()
			} else {
				host = strings.TrimSuffix(names[0], ".")
			}

		} else {
			// We don't store a host for a wide range of IPs
			host = ""
		}
	}

	// Special case: user gave us a range like 10.0.0.0/8 so we don't have a host associated to it (computed just above)
	// If user adds an alias to this access, we will end up trying to ssh an alias that doesn't resolve
	if typeGiven == "IPNET" && host == "" && alias != "" {
		return ba, fmt.Errorf("you cannot add an alias to an IP range")
	}

	ba = &Access{
		Host:   host,
		Prefix: hostIPNetstr,
		Alias:  alias,
		User:   user,
		Port:   intPort,
		IP:     hostIP,
	}

	return
}

// BuildSBAccessFromUserInput deserializes a 'user@host:port' string into a SBAccess struct
func BuildSBAccessFromUserInput(access string) (ba *Access, err error) {

	user, host, port, err := splitUserInput(access, false)
	if err != nil {
		return
	}

	// The alias variable is always empty from user input
	// It is only defined when a user builds an access to store in database
	ba, err = BuildSBAccess(host, user, strconv.Itoa(port), "", false)

	return
}

// GetAllAccesses returns all access in the database
func GetAllAccesses(db *gorm.DB) (accesses []*Access, err error) {

	// Select
	err = db.Find(&accesses).Error

	return
}

// IsAValidSBAccessFromUserInput checks if the provided argument is of the form 'user@host[:port]'
// (or pretty much anything in case of an alias)
func IsAValidSBAccessFromUserInput(access string) bool {
	_, _, _, err := splitUserInput(access, false)
	return err == nil
}

// LoadSBAccess loads a sb access stored in database from user input
func LoadSBAccess(host, user, port string, db *gorm.DB) (ba *Access, err error) {

	// We'll actually call BuildSBAccess with all that to get the host resolution built in that function
	a, err := BuildSBAccess(host, user, port, "", true)
	if err != nil {
		return
	}

	// Instantiate the access where we will load the database result
	ba = &Access{}

	// Select
	err = db.Where(&a).First(&ba).Error

	return
}

// Delete removes the access from the provided database
func (ba *Access) Delete(db *gorm.DB) (err error) {

	// We delete our access
	return db.Delete(ba).Error

}

// Equals compares the properties of two accesses to determine if they're the same access
func (ba *Access) Equals(a *Access) bool {

	if ba.Host != a.Host {
		return false
	}
	if ba.User != a.User {
		return false
	}
	if ba.Port != a.Port {
		return false
	}

	return true
}

// Save saves the access in the provided database
func (ba *Access) Save(db *gorm.DB) (err error) {
	return db.Save(ba).Error
}

// String returns a pretty print display of the access
func (ba *Access) String() string {
	green := color.New(color.FgGreen).SprintFunc()
	return fmt.Sprintf("%s: %-20s | %s: %-20s | %s: %-20s | %s: %-10s | %s: %-5d", green("Prefix"), ba.Prefix, green("Host"), ba.Host, green("Alias"), ba.Alias, green("User"), ba.User, green("Port"), ba.Port)
}

// ShortString returns a pretty print short display of the access
func (ba *Access) ShortString() string {
	green := color.New(color.FgGreen).SprintFunc()
	alias := ""
	port := ""
	host := ba.Prefix
	if ba.Host != "" {
		host = ba.Host
	}
	if ba.Port != 0 {
		port = fmt.Sprintf(":%d", ba.Port)
	}
	if ba.Alias != "" {
		alias = fmt.Sprintf(" (%s)", ba.Alias)
	}
	return green(fmt.Sprintf("%s@%s%s%s", ba.User, host, port, alias))
}

// String returns a pretty print display of the source
func (s *Source) String() (str string) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	if s.Type == "self" {
		str = fmt.Sprintf("- %s", green("personal access"))
	} else {
		str = fmt.Sprintf("- %s from group %s", green("group access"), yellow(s.Group))
	}
	return
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}

func splitUserInput(userInput string, strictHostCheck bool) (user, host string, port int, err error) {

	// We start by determining if we have a specified port
	hostPortParts := strings.Split(userInput, ":")
	if len(hostPortParts) > 1 {
		portInt, errConv := strconv.Atoi(hostPortParts[1])
		if errConv != nil {
			err = fmt.Errorf("port is not a valid integer")
			return
		}
		port = portInt
	}

	// Then we concentrate on the host part (which could be an alias to a host, by the way)
	userHostParts := strings.Split(hostPortParts[0], "@")

	// No @ present in the user input, we consider input to be an alias of a host
	if len(userHostParts) < 2 {

		// Unless we have a strictHostCheck flag
		if strictHostCheck {
			err = fmt.Errorf("unable to parse access from user input %s: no @ separator found", userInput)
			return
		}

		// We set the alias in the host field, and we won't have more info, so let's return
		host = hostPortParts[0]
		return
	}

	user = userHostParts[0]
	host = userHostParts[1]

	return
}
