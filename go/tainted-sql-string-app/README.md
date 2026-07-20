# tainted-sql-string-app

Cùng cách làm với `go/mvc-sqli-app/`, nhưng dựng quanh
`tp_tainted-sql-string.go`. Bản gốc file này **vẫn còn nguyên** ở `go/`
(không xoá/move) - đây là bản copy để gói thành 1 application hoàn chỉnh.

## Cấu trúc (MVC)

- `main.go` - entry point duy nhất.
- `mvc_router.go` - Router (`http.ServeMux`), 6 route.
- `mvc_controller.go` - Controller: nơi request "từ internet" chạm code lần đầu; đồng thời `recover()` panic từ Model (các hàm gốc tự `panic(err)` khi lỗi DB) thành response 500 thay vì sập server.
- `mvc_model.go` - Model mới (`TableModel`), gọi thẳng vào các hàm sinh-handler có sẵn trong `tp_tainted-sql-string.go` (`DeleteHandler`, `SelectHandler`, `SelectHandler2`, `SelectHandler3` - vulnerable; `DeleteHandlerOk`, `SelectHandlerOk` - an toàn, để đối chứng).
- `mvc_view.go` - `RenderServerError`.
- `mvc_simulated_request.go` - biến thể **không router**: `httptest.NewRequest` dựng request rồi gọi thẳng `controller.ServeHTTP`.
- `rules/` - copy `tainted-sql-string.yaml` (taint mode, bắt buộc) + `gosql-sqli.yaml` (pattern-based, cũng khớp các sink này).

## Luồng request -> lỗi (SQLi)

```
internet
  -> Router (mvc_router.go, /table/delete)
  -> Controller (mvc_controller.go: DeleteTableController)
  -> Model (mvc_model.go: TableModel.Delete)
  -> DeleteHandler (tp_tainted-sql-string.go)
       id := req.URL.Query().Get("Id")                 // SOURCE
       db.Exec("DELETE FROM table WHERE Id = " + id)   // SINK
```

3 biến thể sink được giữ nguyên để test đủ dạng pattern của rule taint:
`SelectHandler` (fmt.Sprintf), `SelectHandler2` (strings.Builder),
`SelectHandler3` (+=). Route `-safe` tương ứng dùng bind parameter, dùng để
kiểm tra rule không báo nhầm (false positive).

## Khác biệt nhỏ so với bản gốc trong `go/`

Bản gốc thiếu import (`strconv`, `strings`), thừa import không dùng
(`crypto/tls`, `encoding/json`, `io/ioutil`, `net/url`), dùng biến `err`
chưa khai báo, và 2 dòng không hợp lệ về kiểu (`"..." + id` với `id` là
`int`; so sánh `id != 3` với `id` là `string`). Bản copy ở đây sửa các lỗi
đó để build được như chương trình thật, giữ nguyên 100% ý nghĩa
vulnerable/safe của từng hàm cũng như toàn bộ comment `// ruleid:` /
`// ok:`.

Chưa build/scan thử được trong sandbox này (không có Go toolchain, chưa
tải xong `semgrep` qua proxy) - nên chạy `semgrep --config rules/ .` và
`go build ./...` ở máy thật để xác nhận.
