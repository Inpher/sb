package models

import (
	"fmt"
	"io/ioutil"
	"os"
	osuser "os/user"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inpher/sb/internal/helpers"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestIsMemberOf(t *testing.T) {

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  "",
		},
		Groups: map[string]*Group{},
	}

	user.BuildGroupsMembership([]string{"bg_owners-o", "bg_groupaclk-aclk", "bg_groupgatek-gk", "bg_groupmember"})

	require.Equal(t, false, user.IsACLKeeperOfGroup("groupgatek"), "User is not ACL keeper of group groupgatek")
	require.Equal(t, false, user.IsGateKeeperOfGroup("groupaclk"), "User is not gate keeper of group groupaclk")
	require.Equal(t, false, user.IsMemberOfGroup("groupgatek"), "User is not member of group groupgatek")
	require.Equal(t, false, user.IsOwnerOfGroup("groupgatek"), "User is not owner of group groupgatek")
	require.Equal(t, true, user.IsACLKeeperOfGroup("groupaclk"), "User is ACL keeper of group groupaclk")
	require.Equal(t, true, user.IsGateKeeperOfGroup("groupgatek"), "User is gate keeper of group groupgatek")
	require.Equal(t, true, user.IsMemberOfGroup("groupmember"), "User is not member of group groupmember")
	require.Equal(t, true, user.IsOwnerOfGroup("owners"), "User is not owner of group owners")

	groups, err := user.GetAllGroups()
	require.NoError(t, err, "An unexpected error occurred when calling GetAllGroups()")
	require.Equal(t, user.Groups, groups, "The groups returned are not valid")
}

func TestFilePathes(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{},
	}

	require.Equal(t, fmt.Sprintf("%s/logs.db", user.User.HomeDir), user.GetLocalLogDatabasePath(), "Log database path is not valid")
	require.Equal(t, fmt.Sprintf("%s/accesses.db", user.User.HomeDir), user.getDatabaseAccessFilePath(), "Access database path is not valid")
	require.Equal(t, fmt.Sprintf("%s/.ssh/authorized_keys", user.User.HomeDir), user.getAuthorizedKeysFilePathes(), "Authorized keys path is not valid")
	require.Equal(t, fmt.Sprintf("%s/ttyrecs", user.User.HomeDir), user.GetTtyrecDirectory(), "Ttyrecs directory is not valid")

	expectedKeyFilepathes := []string{
		fmt.Sprintf("%s/test_assets/user/.ssh/id_ed25519", homeDir),
		fmt.Sprintf("%s/test_assets/user/.ssh/id_invalid", homeDir),
		fmt.Sprintf("%s/test_assets/user/.ssh/id_rsa", homeDir),
	}

	keyFiles, err := user.getKeyFilePathes()
	require.NoError(t, err, "An unexpected error occurred when calling getKeyFilePathes")
	require.Equal(t, expectedKeyFilepathes, keyFiles, "Keyfile pathes are not valid")

	expectedPubkeyFilepathes := []string{
		fmt.Sprintf("%s/test_assets/user/.ssh/id_ed25519.pub", homeDir),
		fmt.Sprintf("%s/test_assets/user/.ssh/id_invalid.pub", homeDir),
		fmt.Sprintf("%s/test_assets/user/.ssh/id_rsa.pub", homeDir),
	}

	pubkeyFiles, err := user.getPubKeyFilePathes()
	require.NoError(t, err, "An unexpected error occurred when calling getKeyFilePathes")
	require.Equal(t, expectedPubkeyFilepathes, pubkeyFiles, "Keyfile pathes are not valid")

	err = user.OverrideDatabaseAccessFilePath(":memory:")
	require.NoError(t, err, "An error occurred while calling OverrideDatabaseAccessFilePath")
	require.Equal(t, ":memory:", user.getDatabaseAccessFilePath(), "Overridden access database path is not valid")

	_, _, err = user.DisplayPubKeys("NOT_A_VALID_TYPE")
	require.Error(t, err, fmt.Errorf("unknown type of key"), "NOT_A_VALID_TYPE is not a valid type")
	_, _, err = user.DisplayPubKeys("ingress")
	require.NoError(t, err, "An unexpected error occurred while calling DisplayPubKeys()")
	_, _, err = user.DisplayPubKeys("egress")
	require.NoError(t, err, "An unexpected error occurred while calling DisplayPubKeys()")
}

func TestListPubKeys(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{},
	}

	pubKeys := map[string][]string{
		"ingress": {
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFxu5J1fpfRBHe/2JKreeDGgJlMZji3n97fYm3KJt8Yv sb@localhost",
		},
		"egress": {
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFxu5J1fpfRBHe/2JKreeDGgJlMZji3n97fYm3KJt8Yv sb@localhost",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDDviAgF0HG8m+Fu93Ob0ZgNsboHED1FEi7/LhakVO55Jka0HVV/dKm1Dg+X0+pHlKNteRrLjBT9MA8+cjTdpxCYj/jWovlUcBqZupJTi+xvSGP4q2flZdKTUh+D/bhTwcrQ910BwAzR9iMGqny3m4F62GUTQayhNMHpkOl6wicdwuMN6BYLrcm5qy9tpq0IrBYBWPyi/7knbMNTEH0UqjIIAfrO5ZHlfRs6jJ5R9gMBuJ/C4PIslzIG8WCyzS5kKrSz14xBldcj63eHtoB1ZU6RuaN4OluJLzdFFkRfGsVWQ6sVhpIMAJRCddRD2oACeHzlZiA7k32ddUKuw4Y3v1B sb@localhost",
		},
	}

	_, err := user.listPubKeys("NOT_A_VALID_TYPE")
	require.Error(t, err, fmt.Errorf("unknown type of key"), "NOT_A_VALID_TYPE is not a valid type")

	for keyType, keys := range pubKeys {

		expectedResult := make([]helpers.PublicKey, 0)
		for _, key := range keys {
			publicKey, comment, options, rest, _ := ssh.ParseAuthorizedKey([]byte(key))

			expectedResult = append(expectedResult, helpers.PublicKey{
				PublicKey: publicKey,
				Comment:   comment,
				Options:   options,
				Rest:      rest,
			})
		}

		keys, err := user.listPubKeys(keyType)
		require.NoError(t, err, "An unexpected error occurred when calling listPubKeys")
		require.Equal(t, expectedResult, keys, "The array of helpers.PublicKey is not valid")

	}

}

func TestManageKeys(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{},
	}

	// Start by overriding the authorized_keys file path
	file, _ := ioutil.TempFile("/tmp", "")
	user.OverrideAuthorizedKeysFilePath(file.Name())

	// Add a new ingress key
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDDviAgF0HG8m+Fu93Ob0ZgNsboHED1FEi7/LhakVO55Jka0HVV/dKm1Dg+X0+pHlKNteRrLjBT9MA8+cjTdpxCYj/jWovlUcBqZupJTi+xvSGP4q2flZdKTUh+D/bhTwcrQ910BwAzR9iMGqny3m4F62GUTQayhNMHpkOl6wicdwuMN6BYLrcm5qy9tpq0IrBYBWPyi/7knbMNTEH0UqjIIAfrO5ZHlfRs6jJ5R9gMBuJ/C4PIslzIG8WCyzS5kKrSz14xBldcj63eHtoB1ZU6RuaN4OluJLzdFFkRfGsVWQ6sVhpIMAJRCddRD2oACeHzlZiA7k32ddUKuw4Y3v1B sb@localhost"
	publicKey, comment, options, rest, _ := ssh.ParseAuthorizedKey([]byte(pubKey))
	pubkeyHelper := helpers.PublicKey{
		PublicKey: publicKey,
		Comment:   comment,
		Options:   options,
		Rest:      rest,
	}

	err := user.AddIngressKey(pubkeyHelper.String())
	require.NoError(t, err, "An unexpected error occurred when adding a new ingress public key")

	// List the ingress keys
	keys, err := user.listPubKeys("ingress")
	require.NoError(t, err, "An unexpected error occurred when listing ingress public keys")
	require.Equal(t, []helpers.PublicKey{pubkeyHelper}, keys, "The added ingress key is not present in the final file")

	// Delete the ingress key
	err = user.DeletePubKey("ingress", pubkeyHelper)
	require.NoError(t, err, "An unexpected error occurred when deleting an ingress public key")

	// Relist the ingress keys to check it's disappeared
	keys, err = user.listPubKeys("ingress")
	require.NoError(t, err, "An unexpected error occurred when listing ingress public keys")
	require.Equal(t, 0, len(keys), "The added ingress key is still present in the final file")

	// Delete a wrong type
	err = user.DeletePubKey("NOT_A_VALID_TYPE", pubkeyHelper)
	require.Error(t, err, "NOT_A_VALID_TYPE is not a valid type and should have raise an error")

	// Add an egress pubkey to try the delete method
	egressPubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIC2FqWH9g+71ul7s2TPuXP+GoGt+HapvY+1pWc1f0uVj sb@localhost"
	publicKey, comment, options, rest, _ = ssh.ParseAuthorizedKey([]byte(egressPubKey))
	pubkeyHelper = helpers.PublicKey{
		PublicKey: publicKey,
		Comment:   comment,
		Options:   options,
		Rest:      rest,
	}

	// Delete the ingress key
	err = user.DeletePubKey("egress", pubkeyHelper)
	require.NoError(t, err, "An unexpected error occurred when deleting an ingress public key")
}

func TestShortString(t *testing.T) {

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  "",
		},
		Groups: map[string]*Group{},
	}

	require.Equal(t, user.User.Username, user.ShortString(), "The ShortString() func is not valid")

}

func TestLoadCurrentUser(t *testing.T) {
	user, err := LoadCurrentUser()
	require.NoError(t, err, "An unexpected error occurred while calling LoadCurrentUser()")
	require.Equal(t, os.Getenv("USER"), user.User.Username, "The LoadCurrentUser doesn't match the current user in environment variable")
}

func TestGetAllSBUsers(t *testing.T) {

	// Build a valid path for the subsequent tests
	_, filename, _, _ := runtime.Caller(0)
	etcPasswdPath := fmt.Sprintf("%s/test_assets/user/passwd", filepath.Dir(filename))
	helpers.SetEtcPasswdFilePath(etcPasswdPath)

	users, err := GetAllSBUsers()
	require.NoError(t, err, "An unexpected error occurred when calling GetAllSBUsers")
	require.Equal(t, 5, len(users), "GetAllSBUsers returned an invalid value")

}

func TestInvalidUser(t *testing.T) {
	_, err := LoadUser("INVALID_USERNAME")
	require.Error(t, err, "INVALID_USERNAME is not a real user")
}

func TestLastSSHSessions(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{},
	}

	user.OverrideDatabaseAccessFilePath(":memory:")

	sessions, err := user.GetLastSSHSessions(20)
	require.NoError(t, err, "An unexpected error occurred when calling GetLastSSHSessions()")
	require.Equal(t, 0, len(sessions), "There shouldn't be any session")
}

func TestUserAccesses(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{},
	}

	// Add with classic DB
	user.OverrideDatabaseAccessFilePath(":memory:")

	_, err := user.AddAccess("test.com", "root", "22", "test", "Added for tests")
	require.NoError(t, err, "An unexpected error occurred when calling AddAccess")

	_, err = user.AddAccess("NOT_A_VALID_HOST", "root", "22", "test", "Added for tests")
	require.Error(t, err, fmt.Errorf("host is neither an IP, a prefix or a resolvable host"), "An error should have occurred when calling AddAccess")

	accesses, err := user.GetAccesses()
	require.NoError(t, err, "An unexpected error occurred when calling GetAccesses")
	require.Equal(t, 0, len(accesses[0].Accesses), "There should not be any access available")

	_, err = user.DeleteAccess("test.com", "root", "22")
	require.Error(t, err, fmt.Errorf("record not found"), "Record shouldn't have been found in the memory database")

	// Working with a gorm.DB handler
	db, err := GetAccessGormDB(":memory:")
	require.NoError(t, err, "An unexpected error occurred when getting the database handler")

	ba, err := user.AddAccess("test.com", "root", "22", "test", "Added for tests", db)
	require.NoError(t, err, "An unexpected error occurred when calling AddAccess")
	// We reset the data we don't store in database for comparaison later
	ba.IP = nil

	accesses, err = user.GetAccesses(db)
	require.NoError(t, err, "An unexpected error occurred when calling GetAccesses")
	require.Equal(t, 1, len(accesses[0].Accesses), "There should only be one access available")
	require.Equal(t, ba, accesses[0].Accesses[0], "There should only be one access available")

	_, err = user.DeleteAccess("test.com", "root", "22", db)
	require.NoError(t, err, "An unexpected error occurred when calling DeleteAccess")
}

func TestHasAccesses(t *testing.T) {

	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
		Groups: map[string]*Group{
			"developers": {
				Name:       "developers",
				SystemName: "bg_developers",
				Member:     true,
				Owner:      false,
				ACLKeeper:  false,
				GateKeeper: false,
			},
		},
	}

	user.Groups["developers"].OverrideDatabaseAccessFilePath(":memory:")
	user.Groups["developers"].OverrideKeyFilesRootDir(fmt.Sprintf("%s/.ssh", user.User.HomeDir))

	// Working with a gorm.DB handler
	db, err := GetAccessGormDB(":memory:")
	require.NoError(t, err, "An unexpected error occurred when getting the database handler")

	ba, err := user.AddAccess("test.com", "root", "22", "test", "Added for tests", db)
	require.NoError(t, err, "An unexpected error occurred when calling AddAccess")

	baCf, err := user.AddAccess("one.one.one.one", "root", "22", "one", "Added for tests", db)
	require.NoError(t, err, "An unexpected error occurred when calling AddAccess")

	baCf22022, err := user.AddAccess("one.one.one.one", "test", "22022", "cf", "Added for tests", db)
	require.NoError(t, err, "An unexpected error occurred when calling AddAccess")

	unauthorizedAccessHost, _ := BuildSBAccess("meow.com", "test", "22022", "", false)
	unauthorizedAccessRange, _ := BuildSBAccess("10.0.0.0/8", "test", "22022", "", false)
	unauthorizedAccessAlias, _ := BuildSBAccess("ALIAS", "", "", "", false)
	unauthorizedAccessUser, _ := BuildSBAccess("test.com", "titi", "22", "", false)
	unauthorizedAccessPort, _ := BuildSBAccess("test.com", "root", "22022", "", false)

	accessInfo, err := user.HasAccess(ba, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, true, "Access should be granted")

	accessInfo, err = user.HasAccess(baCf, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, true, "Access should be granted")

	accessInfo, err = user.HasAccess(baCf22022, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, true, "Access should be granted")

	accessInfo, err = user.HasAccess(unauthorizedAccessHost, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, false, "Access should not be granted")

	accessInfo, err = user.HasAccess(unauthorizedAccessRange, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, false, "Access should not be granted")

	accessInfo, err = user.HasAccess(unauthorizedAccessAlias, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, false, "Access should not be granted")

	accessInfo, err = user.HasAccess(unauthorizedAccessUser, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, false, "Access should not be granted")

	accessInfo, err = user.HasAccess(unauthorizedAccessPort, db)
	require.NoError(t, err, "An unexpected error occurred when calling HasAccess")
	require.Equal(t, accessInfo.Authorized, false, "Access should not be granted")

	_, err = user.DeleteAccess("test.com", "root", "22", db)
	require.NoError(t, err, "An unexpected error occurred when calling DeleteAccess")
}

func TestGetSSHKeyPairsInvalidPath(t *testing.T) {

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  "INVALID_PATH",
		},
		Groups: map[string]*Group{},
	}

	_, err := user.GetSSHKeyPairs()
	require.Error(t, err, fmt.Errorf("open INVALID_PATH: no such file or directory"), "An error should be raised if the file path is invalid")
}

func TestTOTP(t *testing.T) {
	// Guess the working directory
	_, filename, _, _ := runtime.Caller(0)
	homeDir := filepath.Dir(filename)

	user := &User{
		User: &osuser.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  fmt.Sprintf("%s/test_assets/user", homeDir),
		},
	}
	testSecret := "randomstring"
	testEmergencyCodes := []string{"10", "11", "12", "13", "14"}

	enabled, _, _ := user.GetTOTP()
	require.Equal(t, false, enabled, "The TOTP for this user should be disabled")

	user.SetTOTPSecret(testSecret, testEmergencyCodes)

	enabled, secret, emergencyCodes := user.GetTOTP()
	require.Equal(t, true, enabled, "The TOTP for this user should be enabled")
	require.Equal(t, testSecret, secret, "The secret value is unexpected")
	require.Equal(t, testEmergencyCodes, emergencyCodes, "The emergency codes are unexpected")

	err := user.RemoveTOTPSecret()
	require.NoError(t, err, "An unexpected error occurred when removing TOTP")

	enabled, _, _ = user.GetTOTP()
	require.Equal(t, false, enabled, "The TOTP for this user should be disabled")
}
