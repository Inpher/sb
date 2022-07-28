package types

import "github.com/spf13/viper"

// StandardError is a wrapper around string to handle the plugin's custom errors
type StandardError string

// Error returns the error as a string
func (e StandardError) Error() string {
	return string(e)
}

const (
	ErrUnknownCommand   StandardError = "unknown command"
	ErrCommandDisabled  StandardError = "command disabled"
	ErrMissingArguments StandardError = "missing arguments"
)

type SetupOptions struct {
	Name                                   string
	Location                               string
	Hostname                               string
	ReplicationEnabled                     bool
	ReplicationPassphrase                  string
	ReplicationQueueType                   string
	ReplicationQueueGooglePubSubProject    string
	ReplicationQueueGooglePubSubTopic      string
	TTYRecsOffloadEnabled                  bool
	TTYRecsOffloadConfigType               string
	TTYRecsOffloadConfigGCSBucket          string
	TTYRecsOffloadConfigGCSObjectsBasePath string
	TTYRecsOffloadConfigS3Bucket           string
	TTYRecsOffloadConfigS3KeysBasePath     string
}

type TTYRecsOffloadingConfig struct {
	Enabled        bool
	StorageType    string
	StorageOptions *viper.Viper
}

type ReplicationQueueConfig struct {
	Enabled      bool
	QueueType    string
	QueueOptions *viper.Viper
}
