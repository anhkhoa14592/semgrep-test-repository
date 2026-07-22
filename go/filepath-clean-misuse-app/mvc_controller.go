package main

// Controller layer.
//
// FileServeController is where a request "from the internet" first
// touches this codebase. The taint source itself (r.FormValue("file")) is
// read inside the existing Model function (handler); this Controller adds
// the one bit of Controller-level policy - restrict to GET - before
// dispatching down into the Model.

import "net/http"

func FileServeController() http.HandlerFunc {
	model := NewFileModel()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RenderMethodNotAllowed(w)
			return
		}

		model.ReadRequestedFile(w, r)
	}
}
