# SimpleIPLocation

簡單的 IP 位置查詢服務

## 功能

查詢當前 IP 的位置

## 使用

```
go run main.go
```

## 指令

| 參數      | 值            | 說明                                                 |
| --------- | ------------- | ---------------------------------------------------- |
| `--host`  | <IP 位址>     | 監聽的位址(預設 0.0.0.0)                             |
| `--port`  | <port number> | 監聽的埠(預設 80)                                    |
| `--log`   | <檔案路徑>    | log file 存放位置 (不指定時不會存放)                 |
| `--proxy` | <true/false>  | 為 true 時會從 "Header X-FORWARDED-FOR" 判斷 IP 位置 |
