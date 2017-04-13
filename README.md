# K2SSアクセスデータ取り込み

## 概要

1. 取得したZIPファイルを解凍
2. ログファイルからデータを抜き出す
3. フィルタを行い、必要なログのみをデータベースへ登録

## 前提

以下がインストールされていること

- Golang

## 実行準備

1. プロジェクトを`$GOPATH/src`直下に配備する。

2. `dep`を取得し、依存ライブラリを取得する。
```sh
go get -u github.com/golang/dep/...
$GOPATH/bin/dep ensure
```

3. 任意ディレクトリに「config.json」を作成する。  
※設定については[Golang用共通ライブラリ](https://github.com/gotoeveryone/golang)を参照

4. ビルドして実行する（Windowsの場合はバイナリ名の末尾に`.exe`が必要）。

```sh
go build -o add-access-detail
add-access-detail --conf=<config.json put directory>
```