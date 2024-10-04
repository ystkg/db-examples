# Goのsql.DBは、いつプールに戻しているのか

## 概要

- Goの標準ライブラリの `database/sql` で中心となる sql.DB ではデフォルトでコネクションプールが使われるようになっていて明示的な操作は不要
- どのタイミングでプールに返却されているのかについてパターンを整理

## Docker

データベースはPostgreSQLのDockerコンテナを使用する

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/docker-compose.yml#L1-L10

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

- *sql.Conn の `Close()` でプールに返却されるパターン。トランザクションありでINSERTを連続して2回実行する例

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0201.go#L11-L59

- 処理の流れを追いやすくするため、関心事だけに絞って、コネクションプールの状態をログ出力させて確認

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0202.go#L9-L24

```shell
go run . ex0202
```

```json
{"time":"2024-09-10T12:12:36.181796186+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-09-10T12:12:36.181937016+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- `conn.Close()` の前後でInUseからIdleに移っている。つまりプールに返却されている

https://github.com/golang/go/blob/go1.23.1/src/database/sql/sql.go#L1193-L1196

## \*sql.DB/*sql.Tx

- CommitもしくはRollbackでプールに返却されるパータン（*sql.Txに `Close()` はない）

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0203.go#L12-L20

```shell
go run . ex0203
```

```json
{"time":"2024-10-03T18:48:28.335019383+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T18:48:28.336647909+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- `tx.Commit()` の後でプールに返却されている

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0204.go#L12-L18

```shell
go run . ex0204
```

```json
{"time":"2024-10-03T18:48:29.981928809+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-10-03T18:48:29.982618874+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- `tx.Rollback()` の後でプールに返却されている

## \*sql.Conn/*sql.Tx

- \*sql.DBでBeginTxした場合と*sql.ConnでBeginTxした場合とで異なる
  - CommitもしくはRollbackではプールに返却されず、*sql.Connの `Close()` するまでは返却されない

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

- 実行毎にプールに返却されるパータン

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0206.go#L12-L18

```shell
go run . ex0206
```

```json
{"time":"2024-09-10T12:13:49.153025371+09:00","level":"INFO","msg":"before","Open":1,"InUse":0,"Idle":1}
{"time":"2024-09-10T12:13:49.156507279+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- InUseが0になっている

## DB.QueryRowContext

- row.Scan()でプールに返却されるパータン

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0207.go#L12-L20

```shell
go run . ex0207
```

```json
{"time":"2024-09-10T12:14:18.114116303+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0}
{"time":"2024-09-10T12:14:18.114284424+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1,"id":1,"name":"shop1"}
```

## DB.QueryContext

- `rows.Next()` が false になったタイミングでプールに返却されるパータン
  - 処理の流れを追いやすくするため、for文を使わずにループを展開

https://github.com/ystkg/db-examples/blob/71ee2b2fcb12ecb81da92a7ff1b9e3f29a4fd427/ex02/ex0208.go#L14-L32

```shell
go run . ex0208
```

```json
{"time":"2024-09-10T12:14:48.961659094+09:00","level":"INFO","msg":"before","Open":1,"InUse":1,"Idle":0,"id":2,"name":"shop2"}
{"time":"2024-09-10T12:14:48.961880508+09:00","level":"INFO","msg":"after ","Open":1,"InUse":0,"Idle":1}
```

- もし仮に false になるまで `rows.Next()` を呼ばなかった場合ですが、そのときは `rows.Close()` のタイミングで返却されることになる

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

- ただし、InUseのコネクションはクローズされない

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

- クローズされるコネクションはIdleのみで、InUseはClose()でプールに戻されることなく、直接クローズされている

## プールの制御

デフォルトでコネクションプールが使われるようになっているが、明示的にプールに戻すタイミングを制御する用途向けとして sql.Conn が用意されいる

https://github.com/golang/go/blob/go1.23.1/src/database/sql/sql.go#L2139-L2146

https://github.com/golang/go/blob/go1.23.1/src/database/sql/sql.go#L1935-L1937

## 関連ドキュメント

<https://go.dev/doc/database/manage-connections>
