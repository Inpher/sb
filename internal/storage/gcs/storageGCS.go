package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type StorageGCS struct {
	bucket       string
	basePath     string
	emulatorHost string
	context      context.Context
	client       *storage.Client
}

func NewStorageGCS(options *viper.Viper) (rs *StorageGCS, err error) {

	bucket := options.GetString("bucket")
	basePath := options.GetString("objects-base-path")
	emulatorHost := options.GetString("emulator-host")

	if bucket == "" {
		err = fmt.Errorf("gcs.bucket option can't be empty")
		return
	}

	rs = &StorageGCS{
		bucket:       bucket,
		basePath:     basePath,
		emulatorHost: emulatorHost,
		context:      context.Background(),
	}

	clientOptions := []option.ClientOption{}
	if rs.emulatorHost != "" {
		rs.emulatorHost = strings.TrimPrefix(rs.emulatorHost, "http://")
		os.Setenv("STORAGE_EMULATOR_HOST", rs.emulatorHost)
		clientOptions = append(clientOptions, option.WithEndpoint(fmt.Sprintf("http://%s/storage/v1/", rs.emulatorHost)))
	}

	rs.client, err = storage.NewClient(rs.context, clientOptions...)
	if err != nil {
		return
	}

	return
}

func (r *StorageGCS) GetFromStorage(object, outputFilePath string) (err error) {
	if r.client == nil {
		return fmt.Errorf("storage GCS hasn't been initialized")
	}

	rc, err := r.client.Bucket(r.bucket).Object(filepath.Join(r.basePath, object)).NewReader(r.context)
	if err != nil {
		return
	}

	file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, rc)
	if err != nil {
		return
	}

	err = rc.Close()
	if err != nil {
		return
	}

	return
}

func (r *StorageGCS) PushToStorage(object, inputFilePath string) (err error) {
	if r.client == nil {
		return fmt.Errorf("storage GCS hasn't been initialized")
	}

	// Open the encrypted file
	file, err := os.Open(inputFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	// Get a writer on GCS
	writer := r.client.Bucket(r.bucket).Object(filepath.Join(r.basePath, object)).NewWriter(r.context)

	// Copy the file to GCS
	_, err = io.Copy(writer, file)
	if err != nil {
		return
	}

	// Close the distant file
	err = writer.Close()
	if err != nil {
		return err
	}

	return
}
