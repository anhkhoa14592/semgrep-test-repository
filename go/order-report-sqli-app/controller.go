package main

// Controller layer: nơi request "từ internet" chạm vào code lần đầu.

import (
	"encoding/json"
	"net/http"
)

func OrderReportController(svc *OrderReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// SOURCE: query param do user kiểm soát hoàn toàn.
		status := r.URL.Query().Get("status")
		if status == "" {
			http.Error(w, "missing status", http.StatusBadRequest)
			return
		}

		orders, err := svc.Report(status)
		if err != nil {
			http.Error(w, "report failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(orders)
	}
}
