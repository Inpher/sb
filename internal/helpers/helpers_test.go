package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testParseArguments struct {
	i []string
	o testParseArgumentsOutputData
}

type testParseArgumentsOutputData struct {
	client          string
	clientArguments []string
	sbArguments     map[string]bool
	arguments       []string
	err             error
}

func TestParseArguments(t *testing.T) {

	tests := []testParseArguments{
		{
			i: []string{}, // This should never happen
			o: testParseArgumentsOutputData{
				client: "ssh",
			},
		},
		{
			i: []string{"sb"},
			o: testParseArgumentsOutputData{
				client: "ssh",
			},
		},
		{
			i: []string{"sb", "self accesses list"},
			o: testParseArgumentsOutputData{
				client:    "ssh",
				arguments: []string{"self accesses list"},
			},
		},
		{
			i: []string{"sb", "-c"},
			o: testParseArgumentsOutputData{
				client: "ssh",
			},
		},
		{
			i: []string{"sb", "-c", "self accesses list 'test"},
			o: testParseArgumentsOutputData{
				err: fmt.Errorf("unclosed quote in command line: self accesses list 'test"),
			},
		},
		{
			i: []string{"sb", "-c", "self accesses list \\'test"},
			o: testParseArgumentsOutputData{
				client:    "ssh",
				arguments: []string{"self accesses list 'test"},
			},
		},
		{
			i: []string{"sb", "-c", "self accesses list"},
			o: testParseArgumentsOutputData{
				client:    "ssh",
				arguments: []string{"self accesses list"},
			},
		},
		{
			i: []string{"sb", "-c", "self ingress-key add --public-key '\"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFxu5J1fpfRBHe/2JKreeDGgJlMZji3n97fYm3KJt8Yv sb@localhost\"'"},
			o: testParseArgumentsOutputData{
				client:    "ssh",
				arguments: []string{"self ingress-key add", "--public-key", "\"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFxu5J1fpfRBHe/2JKreeDGgJlMZji3n97fYm3KJt8Yv sb@localhost\""},
			},
		},
	}

	for i, test := range tests {
		c, ca, ba, a, err := ParseArguments(test.i)
		if test.o.err == nil {
			require.NoError(t, err, fmt.Sprintf("There was an unexpected error parsing arguments on test %d", i+1))
		} else {
			require.EqualError(t, err, test.o.err.Error(), fmt.Sprintf("An error should have been returned on test %d", i+1))
			continue
		}

		require.Equal(t, test.o.client, c, fmt.Sprintf("The client has been badly computed on test %d", i+1))
		require.Equal(t, test.o.clientArguments, ca, fmt.Sprintf("The client arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.sbArguments, ba, fmt.Sprintf("The sb arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.arguments, a, fmt.Sprintf("The remaining arguments were badly computed on test %d", i+1))
	}

}

func TestParseArgumentsMosh(t *testing.T) {

	tests := []testParseArguments{
		{
			i: []string{"sb", "-c", "-v -i selfAddIngressKey"},
			o: testParseArgumentsOutputData{
				client:      "ssh",
				sbArguments: map[string]bool{"verbose": true, "interactive": true},
				arguments:   []string{"selfAddIngressKey"},
			},
		},
	}

	for i, test := range tests {
		c, ca, ba, a, err := ParseArguments(test.i)
		if test.o.err == nil {
			require.NoError(t, err, fmt.Sprintf("There was an unexpected error parsing arguments on test %d", i+1))
		} else {
			require.EqualError(t, err, test.o.err.Error(), fmt.Sprintf("An error should have been returned on test %d", i+1))
			continue
		}

		require.Equal(t, test.o.client, c, fmt.Sprintf("The client has been badly computed on test %d", i+1))
		require.Equal(t, test.o.clientArguments, ca, fmt.Sprintf("The client arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.sbArguments, ba, fmt.Sprintf("The sb arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.arguments, a, fmt.Sprintf("The remaining arguments were badly computed on test %d", i+1))
	}

}

func TestParseArgumentsSBArguments(t *testing.T) {

	tests := []testParseArguments{
		{
			i: []string{"sb", "-c", "mosh-server new -s -v -i 127.0.0.1 -c 256 -l LANG=en_US.UTF-8 -- selfListAccesses"},
			o: testParseArgumentsOutputData{
				client:          "mosh",
				clientArguments: []string{"-s", "-v", "-i", "127.0.0.1", "-c", "256", "-l", "LANG=en_US.UTF-8"},
				arguments:       []string{"selfListAccesses"},
			},
		},
		{
			i: []string{"sb", "-c", "mosh-server 'new' '-s' '-c' '256' '-l' 'LANG=en_US.UTF-8' '--' 'selfListAccesses'"},
			o: testParseArgumentsOutputData{
				client:          "mosh",
				clientArguments: []string{"-s", "-c", "256", "-l", "LANG=en_US.UTF-8"},
				arguments:       []string{"selfListAccesses"},
			},
		},
		{
			i: []string{"sb", "-c", "mosh-server -s -c 256 -l LANG=en_US.UTF-8 -- selfListAccesses"}, // This should never happen
			o: testParseArgumentsOutputData{
				client:          "mosh",
				clientArguments: []string{"-s", "-c", "256", "-l", "LANG=en_US.UTF-8"},
				arguments:       []string{"selfListAccesses"},
			},
		},
		{
			i: []string{"sb", "-c", "mosh-server -s -c 256 -p 22 -l LANG=en_US.UTF-8 -- selfListAccesses"}, // Check if we really dropped the port
			o: testParseArgumentsOutputData{
				client:          "mosh",
				clientArguments: []string{"-s", "-c", "256", "-l", "LANG=en_US.UTF-8"},
				arguments:       []string{"selfListAccesses"},
			},
		},
	}

	for i, test := range tests {
		c, ca, ba, a, err := ParseArguments(test.i)
		if test.o.err == nil {
			require.NoError(t, err, fmt.Sprintf("There was an unexpected error parsing arguments on test %d", i+1))
		} else {
			require.EqualError(t, err, test.o.err.Error(), fmt.Sprintf("An error should have been returned on test %d", i+1))
			continue
		}

		require.Equal(t, test.o.client, c, fmt.Sprintf("The client has been badly computed on test %d", i+1))
		require.Equal(t, test.o.clientArguments, ca, fmt.Sprintf("The client arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.sbArguments, ba, fmt.Sprintf("The sb arguments were badly computed on test %d", i+1))
		require.Equal(t, test.o.arguments, a, fmt.Sprintf("The remaining arguments were badly computed on test %d", i+1))
	}

}
