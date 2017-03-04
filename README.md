# K2SSアクセスデータ取り込み

## 概要

1. 取得したZIPファイルを解凍
2. ログファイルからデータを抜き出す
3. フィルタを行い、必要なログのみをデータベースへ登録

## 前提

以下がインストールされていること

- Golang

## 実行準備

任意ディレクトリに「config.json」を作成する。  
[こちら](https://github.com/gotoeveryone/golang "gotoeveryone/golang") を参照のこと。  

以下コマンドで実行する（Windowsの場合はバイナリ名の末尾に`.exe`が必要）。

```sh
# depの最新を取得
$GOPATH/dep ensure -update
go build -o add-access-detail
add-access-detail
```