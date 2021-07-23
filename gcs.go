package gcsproxy

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

const (
	mkdirPerm = 0755
)

func DownloadGCSObject(dir, filePath string, bucket *storage.BucketHandle) error {
	ctx := context.Background()
	object := bucket.Object(filePath)

	or, err := object.NewReader(ctx)
	if err != nil {
		return err
	}
	defer or.Close()

	fp := filepath.Join(dir, filePath)

	writeDir := filepath.Dir(fp)
	if _, err := os.Stat(writeDir); os.IsNotExist(err) {
		err = os.MkdirAll(writeDir, mkdirPerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, or)
	if err != nil {
		return err
	}

	return nil
}
