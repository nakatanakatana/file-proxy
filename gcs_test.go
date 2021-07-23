package gcsproxy_test

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/storage"
	gcsproxy "github.com/nakatanakatana/gcs-proxy"
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
		err = gcsproxy.DownloadGCSObject("./tmp", "hoge/fuga/vector.toml", gcsBucket)
		if err != nil {
			log.Println(err)
			t.Fail()
		}
	})
}
