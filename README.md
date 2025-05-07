# 個人用単語索引アプリ
個人用の単語索引アプリです．
後から調べようと思った単語の保存と編集を出来るようにします．
webアプリとして実装します．
## 技術選定
- フロントエンド
  - Next.js
- バックエンド
  - Golang

## 起動方法
1. MySQLを立ち上げる
```bash
docker-compose- up -d
```
```bash
docker-compose run  mysql-cli
```
2. goコンパイル＆実行
```bash
go run .
```