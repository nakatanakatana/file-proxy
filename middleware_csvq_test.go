package gcsproxy_test

import (
	"log"
	"testing"

	"github.com/joho/sqltocsv"
	_ "github.com/mithrandie/csvq-driver"
	gcsproxy "github.com/nakatanakatana/gcs-proxy"
)

func TestCSVQ(t *testing.T) {
	t.Skip()
	t.Parallel()

	rows, err := gcsproxy.CSVQ("./tmp", "select * from sample limit 1")
	// rows, err := gcsproxy.CSVQ("./tmp", "select 1")
	if err != nil || rows.Err() != nil {
		t.Fail()
	}
	defer rows.Close()

	str, err := sqltocsv.WriteString(rows)
	if err != nil {
		t.Fail()
	}

	log.Println(str)
}
