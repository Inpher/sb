package storage

import (
	"fmt"

	"github.com/inpher/sb/internal/storage/gcs"
	"github.com/inpher/sb/internal/storage/s3"
	"github.com/inpher/sb/internal/types"
	"github.com/pkg/errors"
)

type Storage interface {
	GetFromStorage(string, string) error
	PushToStorage(string, string) error
}

func GetStorage(config *types.TTYRecsOffloadingConfig) (rs Storage, err error) {

	if !config.Enabled {
		err = fmt.Errorf("TTYRecs offloading is disabled")
		return
	}

	switch config.StorageType {
	case "gcs":
		rs, err = gcs.NewStorageGCS(config.StorageOptions)
	case "s3":
		rs, err = s3.NewStorageS3(config.StorageOptions)
	default:
		err = fmt.Errorf("storage type %s is not implemented", config.StorageType)
	}

	if err != nil {
		err = errors.Wrap(err, "error while initializing ttyrecsoffloading")
	}

	return
}
