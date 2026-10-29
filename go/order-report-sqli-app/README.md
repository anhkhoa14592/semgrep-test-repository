# order-report-sqli-app

App mẫu có **đúng 1 finding SQLi (true positive)**, luồng đi qua 4 file để test pipeline
codegraph → ivern `taint_tree` → fairy.

## Luồng request -> lỗi

```
GET /api/orders/report?status=...              router.go     NewRouter
  -> OrderReportController (closure)           controller.go SOURCE: r.URL.Query().Get("status")
  -> (*OrderReportService).Report               service.go    strings.TrimSpace / strings.ToUpper (KHÔNG chặn SQLi)
  -> (*OrderRepository).ListByStatus            repository.go
  -> queryOrdersByStatus                        repository.go SINK: fmt.Sprintf(... '%s' ...) -> db.Query
```

Payload ví dụ: `status=x' OR '1'='1`.

## Kỳ vọng

- Semgrep (`rules/gosql-sqli.yaml`): 1 finding tại `repository.go`, dòng `db.Query(query)`.
- Ivern (`AST_BACKEND=codegraph`): `sink_kind = sql_format`; 1 nhánh `reaches_external_source`, entry `http_route "/api/orders/report"`,
  source `http_request` (`r.URL.Query().Get("status")`), đường đi `main → StartServer → NewRouter → OrderReportController
  → Report → ListByStatus → queryOrdersByStatus`, không validator (TrimSpace/ToUpper không phải validator), guard `status == ""`.
- Fairy: true positive — validator chỉ chuẩn hoá chữ, không loại `'`.
