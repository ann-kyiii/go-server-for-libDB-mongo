# go-server-for-libDB-mongo

図書貸出アプリのサーバプログラムです．  
`golang，echo` で作られています．

このリポジトリは[go-server-for-libDB](https://github.com/ann-kyiii/go-server-for-libDB)のDBをfirestoreからmongoDBを使う形で書き換えたものです．

## サーバの起動方法
サーバは，以下のコマンドでクイック起動できます．
```
go run main_server.go
```

また，ビルドとサーバのバックグラウンド実行は以下のコマンドで行います．
```
go build main_server.go
./main_server.go &
```

バックグラウンド実行したサーバを終了するには，実行中のプロセスを以下のコマンドで探して，2項目のPIDを指定して終了します．
```
ps aux | grep main_server
kill [PID(.main_serverプロセスの2項目の数字))]
```

###  サーバへのアクセス例
APIサーバへのアクセス方法は，`test_command` を例とします．

例：
```
curl [サーバのグローバルIPアドレス]:1313 
```

## mongoDBのセットアップ
あらかじめmongoDBをインストールし，データのcsvファイルを準備します．  
この際，csvファイルのヘッダは以下の通り型を付与したものを使用します．
```
id.int64(),bookName.string(),genre.string(),subGenre.string(),ISBN.string(),find.int64(),sum.int64(),author.string(),publisher.string(),pubdate.string(),exist.string(),locateAt4F.boolean(),withDisc.string(),other.string(),imgURL.string()
```

そして，以下のコマンドでcsvファイルからDBにデータをインポートします．
```bash
mongoimport --db=<YOUR DB NAME> --collection=<YOUR COLLECTION NAME>  --type csv --file <CSV_FILEPATH> --headerline --columnsHaveTypes
```

最後にborrowerの空配列を追加するため，add_borrower.goを実行する．
```
go run add_borrower/add_borrower.go
```

# Docker環境のセットアップ
本リポジトリには`.devcontainer`と`docker`の2つの設定があり，前者が開発用で後者が運用用です．

`.devcontainer`はVSCodeのDev Containerで使うための設定であり，本リポジトリをマウントします．

`docker`はサーバを起動するように設定しているため，`docker compose up`でサーバが起動します．

# データの更新手順
1. 更新データを用意し`initdb/update/updatedata.csv`に配置
  * id列は今のところ手作業で振る
2. `docker exec -it library-app-server-mongodb-1 sh /docker-entrypoint-initdb.d/update/update.sh`
  * コンテナ名が異なる場合は適切なコンテナ名に変更する
