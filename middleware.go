package gcsproxy

import (
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

func GetGCSFile(dir string, bucket *storage.BucketHandle, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "" {
			filePath := path[1:]
			err := DownloadGCSObject(dir, filePath, bucket)
			if err != nil {
				log.Println("DownloadObjectError", dir, filePath, err)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func CreateFileServer(rootPath string) http.Handler {
	return http.FileServer(http.Dir(rootPath))
}
