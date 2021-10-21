package gcsproxy

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/sqltocsv"

	// database-driver.
	_ "github.com/mithrandie/csvq-driver"
)

func CSVQ(repository string, query string) (*sql.Rows, error) {
	reposPath, err := filepath.Abs(repository)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("csvq", reposPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	return db.Query(query)
}

func CSVQFilter(dir string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("csvq")
		if q == "" {
			h.ServeHTTP(w, r)

			return
		}

		path := r.URL.Path
		if path == "" {
			h.ServeHTTP(w, r)

			return
		}

		filePath := filepath.Join(dir, path)
		if _, err := os.Stat(filePath); err != nil {
			h.ServeHTTP(w, r)

			return
		}

		paths := strings.Split(r.URL.Path, "/")
		tmpDir := filepath.Join(dir, filepath.Join(paths[:len(paths)-1]...))

		rows, err := CSVQ(tmpDir, q)
		if err != nil || rows.Err() != nil {
			h.ServeHTTP(w, r)

			return
		}

		defer rows.Close()

		w.Header().Set("Content-type", "text/csv")
		err = sqltocsv.Write(w, rows)
		if err != nil {
			return
		}
	})
}
