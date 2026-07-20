# kafka-order-sqli-app

Cùng phong cách MVC với 3 app trước trong `go/`, nhưng nguồn dữ liệu không
phải HTTP request mà là **Kafka message** (dùng `github.com/segmentio/kafka-go`).
Đây là app hoàn toàn mới (không copy từ file có sẵn trong `go/`).

## Cấu trúc (MVC)

- `main.go` - entry point: mở DB, mở `kafka.Reader` (topic `order-events`), gọi `RunWithRouter`.
- `mvc_router.go` - Router: `TopicRouter` map `topic -> handler` (tương đương router HTTP nhưng route theo topic thay vì path); `RunWithRouter` là vòng lặp `reader.ReadMessage` liên tục - tương đương `ListenAndServe`.
- `mvc_controller.go` - Controller: nơi message "từ Kafka" chạm code lần đầu (`NewOrderEventController`), đọc thẳng `msg.Key`/`msg.Value`. Kèm `NewOrderEventSafeController` để đối chứng.
- `mvc_model.go` - Model: `OrderModel.UpdateOrderStatus` (vulnerable, `fmt.Sprintf` + `db.Exec`) và `UpdateOrderStatusSafe` (bind parameter, an toàn).
- `mvc_view.go` - View: log kết quả xử lý (Kafka consumer không có HTTP response để render).
- `mvc_simulated_message.go` - biến thể **không router**: tự dựng 1 `kafka.Message` (Key/Value) rồi gọi thẳng `controller(msg)`, không cần router, không cần kết nối broker thật.
- `rules/kafka-tainted-sql-string.yaml` - rule semgrep taint mode mới, source đổi từ `*http.Request` sang `kafka.Message`.
- `go.mod` - tham khảo; chưa build-verify (không có Go toolchain / không tải được module qua proxy trong sandbox này).

## Luồng message -> lỗi (SQLi)

```
Kafka broker (topic: order-events)
  -> Router (mvc_router.go: TopicRouter.Dispatch, route theo msg.Topic)
  -> Controller (mvc_controller.go: NewOrderEventController)
       orderID := string(msg.Key)     // SOURCE
       status  := string(msg.Value)   // SOURCE
  -> Model (mvc_model.go: OrderModel.UpdateOrderStatus)
       query := fmt.Sprintf("UPDATE orders SET status = '%s' WHERE order_id = '%s'", status, orderID)
       m.DB.Exec(query)               // SINK
```

Rule `kafka-tainted-sql-string.yaml` (taint mode, `interfile: true`) định
nghĩa source là bất kỳ truy cập `.Value`/`.Key`/`.Headers` trên biến kiểu
`kafka.Message`, và tái dùng đúng sink pattern (fmt.Sprintf / nối chuỗi /
strings.Builder + regex từ khoá SQL) như `tainted-sql-string.yaml` đã có
trong `go/`. Vì source (Controller) và sink (Model) nằm 2 file khác nhau,
đây tiếp tục là ca test cho taint xuyên file / call-graph, chỉ khác nguồn
input.

`NewOrderEventSafeController` + `UpdateOrderStatusSafe` là route đối chứng
(dùng bind parameter) để kiểm tra rule không báo nhầm.
