package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	fproxy "github.com/nakatanakatana/file-proxy"
)

const (
	HTTPReadTimeout  = 300 * time.Second
	HTTPWriteTimeout = 300 * time.Second
)

func main() {
	ctx := context.Background()

	targetDir := os.Getenv("GCS_PROXY_DIR")
	if targetDir == "" {
		log.Fatal("error get GCS_PROXY_DIR")
	}

	targetBucket := os.Getenv("GCS_PROXY_BUCKET")
	if targetBucket == "" {
		log.Fatal("error get GCS_PROXY_BUCKET")
	}

	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	bucket := gcsClient.Bucket(targetBucket)

	mux := http.NewServeMux()
	mux.Handle("/",
		fproxy.GetGCSFile(targetDir, bucket,
			fproxy.CSVQFilter(targetDir,
				fproxy.CreateFileServer(targetDir)),
		),
	)

	svr := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  HTTPReadTimeout,
		WriteTimeout: HTTPWriteTimeout,
	}

	log.Fatal(svr.ListenAndServe())
}
