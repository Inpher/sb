package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	require.Equal(t, "sb", GetSBName(), "The default value of sb name is wrong")
	require.Equal(t, "sb.domain.tld", GetSBHostname(), "The default value of sb hostname is wrong")
	require.Equal(t, "40000:49999", GetMOSHPortsRange(), "The default value of the mosh ports rangeis wrong")
	require.Equal(t, []string{"USER"}, GetEnvironmentVarsToForward(), "The default SSH environmen variables to forward is wrong")
	require.Equal(t, "ttyrec", GetSSHCommand(), "The default value of sb hostname is wrong")
	require.Equal(t, "/opt/sb/sb", GetBinaryPath(), "The default value of the binary path is wrong")
}
