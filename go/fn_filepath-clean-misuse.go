package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {

	// SOURCE
	userPath := r.FormValue("file")

	/*
		Manual mitigation:
		- normalize path
		- reject traversal after clean
		- enforce base directory
	*/

	cleaned := filepath.Clean(userPath)

	/*
		Semgrep still alerts because:
		- source -> filepath.Clean(...)
		- sanitizer only matches:
		      "/" + ...
		- this usage still matches sink directly
	*/

	// reject traversal attempts
	if strings.Contains(cleaned, "..") {
		http.Error(w, "invalid path", 400)
		return
	}

	baseDir := "/var/app/uploads"

	finalPath := filepath.Join(baseDir, cleaned)

	// ensure path stays inside baseDir
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	absFinal, err := filepath.Abs(finalPath)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	if !strings.HasPrefix(absFinal, absBase+string(os.PathSeparator)) &&
		absFinal != absBase {

		http.Error(w, "path escape detected", 400)
		return
	}

	data, err := os.ReadFile(absFinal)
	if err != nil {
		http.Error(w, "read error", 500)
		return
	}

	w.Write(data)
}

func main() {}