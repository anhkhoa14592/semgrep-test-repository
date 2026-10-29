package main

// Service layer: "chuẩn hoá" input trước khi xuống repository.
// TrimSpace/ToUpper KHÔNG chặn được SQLi (không loại dấu ') — validator giả,
// dùng để kiểm tra fairy/ivern có phân biệt được validator thật hay không.

import "strings"

type OrderReportService struct {
	repo *OrderRepository
}

func NewOrderReportService(repo *OrderRepository) *OrderReportService {
	return &OrderReportService{repo: repo}
}

func (s *OrderReportService) Report(status string) ([]Order, error) {
	normalized := strings.ToUpper(strings.TrimSpace(status))
	return s.repo.ListByStatus(normalized)
}
