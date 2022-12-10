package gcsproxy

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

const (
	mkdirPerm = 0755
)

func DownloadGCSObject(dir, filePath string, bucket *storage.BucketHandle) error {
	var (
		objectReader *storage.Reader
		err          error
	)
	ctx := context.Background()

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
				return err
			}
		} else {
			return err
		}
	}
	defer objectReader.Close()

	fp := filepath.Join(dir, targetPath)

	writeDir := filepath.Dir(fp)
	if _, err := os.Stat(writeDir); os.IsNotExist(err) {
		err = os.MkdirAll(writeDir, mkdirPerm)
		if err != nil {
			log.Println("MkdirAll")
			return err
		}
	}

	f, err := os.Create(fp)
	if err != nil {
		log.Println("Create")
		return err
	}

	_, err = io.Copy(f, objectReader)
	if err != nil {
		log.Println("Copy")
		return err
	}

	return nil
}

func GetGCSFile(dir string, bucket *storage.BucketHandle, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "" {
			filePath := path[1:]
			if info, err := os.Stat(filepath.Join(dir, filePath)); err != nil || info.IsDir() {
				err := DownloadGCSObject(dir, filePath, bucket)
				if err != nil {
					log.Println("DownloadObjectError", dir, filePath, err)
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}
