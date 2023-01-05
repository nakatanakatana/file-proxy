package fproxy_test

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/storage"
	fproxy "github.com/nakatanakatana/file-proxy"
)

func TestDownloadGCSObject(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctx := context.Background()

	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	t.Run("exists path", func(t *testing.T) {
		t.Parallel()
		gcsBucket := gcsClient.Bucket("vector_config")
		ctx := context.Background()

		err = fproxy.DownloadGCSObject(ctx, "./tmp", "hoge/fuga/vector.toml", gcsBucket)
		if err != nil {
			log.Println(err)
			t.Fail()
		}
	})
}
