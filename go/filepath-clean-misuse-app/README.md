# filepath-clean-misuse-app

Cùng cách làm với `go/mvc-sqli-app/` và `go/tainted-sql-string-app/`, dựng
quanh `fn_filepath-clean-misuse.go`. Bản gốc file này **vẫn còn nguyên** ở
`go/` (không xoá/move) - đây là bản copy để gói thành 1 application hoàn
chỉnh. Không như 2 app trước, file gốc lần này không có lỗi cú pháp nào -
chỉ bỏ `func main() {}` (đưa vào `main.go` riêng của app).

## Cấu trúc (MVC)

- `main.go` - entry point duy nhất.
- `mvc_router.go` - Router (`http.ServeMux`), 1 route.
- `mvc_controller.go` - Controller: nơi request "từ internet" chạm code lần đầu; áp policy duy nhất ở tầng này (chỉ cho phép GET) trước khi xuống Model.
- `mvc_model.go` - Model mới (`FileModel`, không có state), gọi thẳng vào `handler` có sẵn trong `fn_filepath-clean-misuse.go`.
- `mvc_view.go` - `RenderMethodNotAllowed`.
- `mvc_simulated_request.go` - biến thể **không router**: `httptest.NewRequest` dựng request rồi gọi thẳng `controller.ServeHTTP`.
- `rules/filepath-clean-misuse.yaml` - rule taint mode (`interfile: true`) tương ứng.

## Luồng request -> lỗi (Path Traversal)

```
internet
  -> Router (mvc_router.go, /files/read)
  -> Controller (mvc_controller.go: FileServeController, chỉ cho GET)
  -> Model (mvc_model.go: FileModel.ReadRequestedFile)
  -> handler (fn_filepath-clean-misuse.go)
       userPath := r.FormValue("file")           // SOURCE
       cleaned := filepath.Clean(userPath)       // SINK
       ... (mitigation thủ công không được rule công nhận là sanitizer)
       os.ReadFile(absFinal)
```

Khác với 2 app trước, source và sink ở đây nằm cùng trong `handler` (không
xuyên file) - phần MVC mới thêm vào chủ yếu để test khả năng truy vết
call-graph qua nhiều lớp gọi gián tiếp (Router -> closure Controller ->
method Model -> hàm gốc), chứ không phải taint xuyên file như 2 app SQLi
trước.

Chưa build/scan thử được trong sandbox này (không có Go toolchain, chưa
tải xong `semgrep` qua proxy) - nên chạy `semgrep --config rules/ .` và
`go build ./...` ở máy thật để xác nhận.
