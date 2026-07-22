# mvc-sqli-app

Self-contained sample application, gộp lại từ các file trước đó nằm rải rác
trong `go/`. Bản gốc của `tn_gosql-sqli.go` và `tn_open-redirect.go` **vẫn
còn nguyên** ở `go/` (không xoá/move) - đây là bản copy để gói thành 1
application hoàn chỉnh phục vụ test semgrep + call-graph.

## Cấu trúc (MVC)

- `main.go` - entry point duy nhất của package này (mở DB, gọi `StartServer`).
- `router/` &rarr; `mvc_router.go` - Router (`http.ServeMux`), map route tới controller.
- Controller &rarr; `mvc_controller.go` - nơi request "từ internet" chạm vào code lần đầu (đọc `r.URL.Query().Get("email")`, `r.FormValue(...)`, hoặc forward nguyên `*http.Request`).
- Model &rarr; `mvc_model.go` (`UserModel`, type mới) + `tn_gosql-sqli.go` (`bad2/bad3/bad4/bad5` - hàm sink có sẵn, không sửa).
- View &rarr; `mvc_view.go` - render response.
- `mvc_simulated_request.go` - biến thể **không router**: dùng `httptest.NewRequest` dựng request rồi gọi thẳng `controller.ServeHTTP`.
- `rules/` - copy 3 rule semgrep liên quan trực tiếp tới app này (`gosql-sqli.yaml`, `tainted-sql-string.yaml` - taint mode, `open-redirect.yaml` - taint mode), để scan độc lập không cần trỏ ra `go/`.
- `go.mod` - để tham khảo cấu trúc module; **chưa build-verify được** trong sandbox này (không có Go toolchain, và không tải được `semgrep` qua proxy để chạy thử).

## Luồng request -> lỗi (SQLi)

```
internet
  -> Router (mvc_router.go, /api/users/search)
  -> Controller (mvc_controller.go: SearchUserByEmailController)
       email := r.URL.Query().Get("email")      // SOURCE
  -> Model (mvc_model.go: UserModel.SearchByEmail)
  -> bad3 (tn_gosql-sqli.go)
       query := fmt.Sprintf("... WHERE email='%s'", email)
       db.Exec(query)                            // SINK
```

Rule `tainted-sql-string.yaml` (taint mode, `interfile: true`) khớp cả
source lẫn sink; do source và sink nằm ở 2 file khác nhau (controller vs
model), đây là ca kiểm thử cho khả năng truy vết taint xuyên file/call-graph.

## Khác biệt nhỏ so với bản gốc trong `go/`

Bản gốc `tn_gosql-sqli.go` có vài lỗi cú pháp không ảnh hưởng semgrep nhưng
khiến file không compile được (`query = ...` thiếu `:=`, import `strings`
không dùng tới). Bản copy trong thư mục này đã sửa các lỗi đó để cả app có
thể build được như một chương trình thật; các hàm `bad1-bad5`, `ok1-ok6`
logic giữ nguyên 100%.
