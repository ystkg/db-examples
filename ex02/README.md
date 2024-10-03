# Goのsql.DBは、いつプールに戻しているのか

## 概要

- Goの標準ライブラリの `database/sql` で中心となる sql.DB ではデフォルトでコネクションプールが使われるようになっていて明示的な操作は不要
- プールに返却されるタイミングを整理

## Docker

- データベースはPostgreSQLのDockerコンテナを使用

### データベースのコンテナ起動

```shell
docker-compose up -d
```

もしくはDocker Composeのプラグイン版なら

```shell
docker compose up -d
```

### データベースのコンテナ削除

```shell
docker-compose down
```

もしくはDocker Composeのプラグイン版なら

```shell
docker compose down
```

## テーブル

- 実行時のセットアップ処理で初期化
- 1テーブル（shop）のみ

```mermaid
erDiagram
    shop {
        int id PK
        string name
        datetime created_at
    }
```

## サンプルコードの実行

```shell
go run . サンプル名
```

- サンプル名は大文字小文字の区別なし

例

```shell
go run . ex0201
```

## *sql.Conn

- *sql.Conn の `Close()` でプールに返却
- トランザクションありでINSERTを連続して2回実行
- 実装の全体

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0201.go#L11-L59

- コネクションプールの状態をログ出力させて確認
- 処理の流れを追いやすくするため、関心事だけに絞る

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0202.go#L9-L24

確認するプールの各ステータスです。

https://github.com/golang/go/blob/ed07b321aef7632f956ce991dd10fdd7e1abd827/src/database/sql/sql.go#L1193-L1196

```shell
go run . ex0202
```

```json
{"time":"2024-09-10T12:12:36.181796186+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-09-10T12:12:36.181937016+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- `conn.Close()` の前後でInUseからIdleに移っている。つまりプールに返却されている

## \*sql.DB/*sql.Tx

- CommitもしくはRollbackでプールに返却される。*sql.Txに `Close()` はない

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0203.go#L12-L20

```shell
go run . ex0203
```

```json
{"time":"2024-10-03T18:48:28.335019383+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T18:48:28.336647909+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0204.go#L12-L18

```shell
go run . ex0204
```

```json
{"time":"2024-10-03T18:48:29.981928809+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T18:48:29.982618874+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

## \*sql.Conn/*sql.Tx

- CommitもしくはRollbackでプールに返却されない。*sql.Connの `Close()` で返却
- \*sql.DBでBeginTxした場合と*sql.ConnでBeginTxした場合とで異なる

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0205.go#L12-L23

```shell
go run . ex0205
```

```json
{"time":"2024-10-03T18:52:32.601262941+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T18:52:32.602614139+09:00","level":"INFO","msg":"after ","Open":1,"InUse":1,"Idle":0,"err":null}
```

- `InUse` に残ったままで返却されていない
- エラーも発生していない

## DB.ExecContext

- 実行毎にプールに返却

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0206.go#L12-L18

```shell
go run . ex0206
```

```json
{"time":"2024-09-10T12:13:49.153025371+09:00","level":"INFO","msg":"before","Open":1,"InUse":0,"Idle":1}
{"time":"2024-09-10T12:13:49.156507279+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

## DB.QueryRowContext

- row.Scan()でプールに返却

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0207.go#L12-L20

```shell
go run . ex0207
```

```json
{"time":"2024-09-10T12:14:18.114116303+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-09-10T12:14:18.114284424+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1,"id":1,"name":"shop1"}
```

## DB.QueryContext

- `rows.Next()` が false になったタイミングでプールに返却
- 処理の流れを追いやすくするため、for文を使わずにループを展開

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0208.go#L14-L32

```shell
go run . ex0208
```

```json
{"time":"2024-09-10T12:14:48.961659094+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0,"id":2,"name":"shop2"}
{"time":"2024-09-10T12:14:48.961880508+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- もし仮に false になるまで `rows.Next()` を呼ばなかった場合ですが、そのときは `rows.Close()` のタイミングで返却されました。

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0209.go#L14-L26

```shell
go run . ex0209
```

```json
{"time":"2024-09-10T12:15:09.62954833+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0,"id":1,"name":"shop1"}
{"time":"2024-09-10T12:15:09.62968931+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

## DB.Close

- DB.Closeはプールにあるコネクションをクローズする

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0210.go#L9-L21

```shell
go run . ex0210
```

```json
{"time":"2024-10-03T19:09:47.607473683+09:00","level":"INFO","msg":"before","Open":5,"InUse":0,"Idle":5}
{"time":"2024-10-03T19:09:47.607883584+09:00","level":"INFO","msg":"after ","Open":0,"InUse":0,"Idle":0}
```

- クローズされるコネクションはIdleのみで、InUseはClose()でプールに戻されることなく、直接クローズされる

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0211.go#L9-L35

```shell
go run . ex0211
```

```json
{"time":"2024-10-03T19:33:07.299576312+09:00","level":"INFO","msg":"before","Open":5,"InUse":3,"Idle":2}
{"time":"2024-10-03T19:33:07.299771464+09:00","level":"INFO","msg":"after ","Open":3,"InUse":3,"Idle":0}
{"time":"2024-10-03T19:33:07.299887676+09:00","level":"INFO","msg":"conn 0","Open":2,"InUse":2,"Idle":0}
{"time":"2024-10-03T19:33:07.300011412+09:00","level":"INFO","msg":"conn 1","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T19:33:07.300026982+09:00","level":"INFO","msg":"conn 2","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T19:33:07.300033695+09:00","level":"INFO","msg":"conn 3","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T19:33:07.300089802+09:00","level":"INFO","msg":"conn 4","Open":0,"InUse":0,"Idle":0}
```

## 関連ドキュメント

<https://go.dev/doc/database/manage-connections>
