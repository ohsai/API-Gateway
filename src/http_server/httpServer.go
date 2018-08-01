package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		panic("Error : need one and only one argument specifying port number")
	}
	createServer(os.Args[1:2][0])
}
func createServer(PORT string) {
	http.Handle("/", new(staticHandler))
	http.ListenAndServe(":"+PORT, nil)
}

type staticHandler struct {
	http.Handler
}

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	localPath := "/home/ubuntu/Workspace/haskell/PersonalUse/src/http_server/http_root" + req.URL.Path
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		fmt.Println("HTTPServer error : ", err.Error())
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)
	w.Write(content)
}

func getContentType(localPath string) string {
	var contentType string
	ext := filepath.Ext(localPath)

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}

	return contentType
}
