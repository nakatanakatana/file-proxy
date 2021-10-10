package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	gcsproxy "github.com/nakatanakatana/gcs-proxy"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const (
	cacheCapacity  = 100
	cacheClientTTL = 1 * time.Minute
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

	mem, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cacheCapacity),
	)
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(mem),
		cache.ClientWithTTL(cacheClientTTL),
	)
	if err != nil {
		log.Fatal(err)
	}

	bucket := gcsClient.Bucket(targetBucket)

	mux := http.NewServeMux()
	mux.Handle("/", cacheClient.Middleware(
		gcsproxy.GetGCSFile(targetDir, bucket, gcsproxy.CSVQFilter(targetDir, gcsproxy.CreateFileServer(targetDir)))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
