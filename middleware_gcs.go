package fproxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

const mkdirPerm = 0o755

func DownloadGCSObject(ctx context.Context, dir, filePath string, bucket *storage.BucketHandle) error {
	var (
		objectReader *storage.Reader
		err          error
	)

	targetPath := filePath
	object := bucket.Object(targetPath)

	objectReader, err = object.NewReader(ctx)
	if err != nil {
		if strings.HasSuffix(targetPath, "/") {
			log.Println("fallback to index.html")

			targetPath += "index.html"
			indexObject := bucket.Object(targetPath)

			objectReader, err = indexObject.NewReader(ctx)
			if err != nil {
				log.Println("fallback error")

				return fmt.Errorf("indexObject.NewReader: %w", err)
			}
		} else {
			return fmt.Errorf("object.NewReader: %w", err)
		}
	}
	defer objectReader.Close()

	file := filepath.Join(dir, targetPath)

	writeDir := filepath.Dir(file)
	if _, err := os.Stat(writeDir); os.IsNotExist(err) {
		err = os.MkdirAll(writeDir, mkdirPerm)
		if err != nil {
			return fmt.Errorf("os.MkdirAll: %w", err)
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	_, err = io.Copy(f, objectReader)
	if err != nil {
		return fmt.Errorf("io.Copy; %w", err)
	}

	return nil
}

func GetGCSFile(dir string, bucket *storage.BucketHandle, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "" {
			filePath := path[1:]
			if info, err := os.Stat(filepath.Join(dir, filePath)); err != nil || info.IsDir() {
				err := DownloadGCSObject(r.Context(), dir, filePath, bucket)
				if err != nil {
					log.Println("DownloadObjectError", dir, filePath, err)
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}
