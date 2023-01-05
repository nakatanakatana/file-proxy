package fproxy_test

import (
	"log"
	"testing"

	"github.com/joho/sqltocsv"
	_ "github.com/mithrandie/csvq-driver"
	fproxy "github.com/nakatanakatana/file-proxy"
)

func TestCSVQ(t *testing.T) {
	t.Skip()
	t.Parallel()

	rows, err := fproxy.CSVQ("./tmp", "select * from sample limit 1")
	// rows, err := fproxy.CSVQ("./tmp", "select 1")
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
