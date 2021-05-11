# SimpleIPLocation

簡單的 IP 位置查詢服務

IP 庫由 [DB-IP](https://db-ip.com/) 提供

## 功能

查詢當前 IP 的資訊

## 建置

下載原始碼

```
git clone https://github.com/orangesobeautiful/SimpleIPLocation.git
cd SimpleIPLocation
```

###後端

```
cd backend
go build main.go
```

#### 參數

| 參數       | 值         | 說明                                 |
| ---------- | ---------- | ------------------------------------ |
| `--host`   | <IP 位址>  | 監聽的位址(預設 0.0.0.0)             |
| `--port`   | <連接埠>   | 監聽的埠(預設 80)                    |
| `--log`    | <檔案路徑> | log file 存放位置 (不指定時不會存放) |
| `--stdout` | 有/無      | 使用 stdout，(預設 沒有)             |
| `--debug`  | 有/無      | 使用 debug mode，(預設 無)           |

###前端

```
cd ../frontend
yarn
quasar build
```

#### 參數

| 參數       | 值         | 說明                                 |
| ---------- | ---------- | ------------------------------------ |
| `--host`   | <IP 位址>  | 監聽的位址(預設 0.0.0.0)             |
| `--port`   | <連接埠>   | 監聽的埠(預設 80)                    |
| `--log`    | <檔案路徑> | log file 存放位置 (不指定時不會存放) |
| `--stdout` | 有/無      | 使用 stdout，(預設 沒有)             |
| `--debug`  | 有/無      | 使用 debug mode，(預設 無)           |
