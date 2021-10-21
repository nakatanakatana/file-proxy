package gcsproxy

import (
	"net/http"
)

func CreateFileServer(rootPath string) http.Handler {
	return http.FileServer(http.Dir(rootPath))
}
