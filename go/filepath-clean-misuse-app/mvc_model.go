package main

// Model layer.
//
// FileModel has no state of its own - the pre-existing handler function
// (fn_filepath-clean-misuse.go) already contains the path-building and
// file-read logic, including its insufficient path-traversal mitigation
// (filepath.Clean is not a sanitizer against traversal). This type just
// gives the Controller layer a proper Model entry point to call through.

import "net/http"

type FileModel struct{}

func NewFileModel() *FileModel {
	return &FileModel{}
}

// ReadRequestedFile -> handler: reads "file" straight off the request
// (inside handler itself) and serves it from disk after the existing,
// insufficient mitigation attempt.
func (m *FileModel) ReadRequestedFile(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}
