package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var dir string
var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Port on which the server should listen")
	flag.StringVar(&dir, "cwd", ".", "Root directory of files.")
	flag.Parse()
	directory, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err.Error())
	}
	dir = directory
	fmt.Printf("Serving files from %s.\n", dir)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveFile)
	portStr := fmt.Sprintf(":%v", port)
	fmt.Printf("Running on port %s\n", portStr)
	err := http.ListenAndServe(portStr, mux)
	log.Fatal(err)
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	fullPath := filepath.Join(dir, r.URL.Path)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}
	f, _ := os.Open(fullPath)
	w.Header().Set("Content-Type", getFileContentType(f))
	f.Close()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(fullPath)))
	log.Printf("Serving file %s", fullPath)
	http.ServeFile(w, r, fullPath)
	return

}

func getFileContentType(out *os.File) string {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "application/octet-stream"
	}

	return http.DetectContentType(buffer)
}
