# 概要

edinetからデータを取ってきて従業員情報を表示するCLI  
取ってきたデータをsqliteに保管する。

# ビルド方法

go build -o edinet

# 使い方

## 環境変数の設定

EDINETにアクセスするためのtokenを環境変数にセットする

```bash
export EDINET_API_TOKEN=xxxxxx
```

<https://disclosure2dl.edinet-fsa.go.jp/guide/static/disclosure/WZEK0110.html>

v2のAPIを使っている。

## ビルド方法

まずはUIのビルド

```bash
cd ui/edinet_search_ui
npm run build
```

その後バイナリ生成

```bash
go  generate
go build -o edinet
```

## DB作成

まずはedinetから取ってきたデータでdbを作るために以下のコマンドを一度実行  
startからendまでの日付の有価証券データを取ってきてDBにいれる  
ordinanceCodeが010、form_codeが030000のデータのみをまずは取るようにしている  

```bash
edinet createdb -s "2023-12-12" -e "2023-12-14"
```

このコマンドは実行されるたびにsqliteのDBを削除して一から作るので、一度実行したら新しいデータを入れたいとき以外は動かさない。  
時間があれば更新されるようにもしたい  

DBファイルは実行した場所に作られる

## 　検索

上記で作ったDBファイルが有る場所で実行  

こんな感じで会社名で検索できる。あいまい検索なのである程度は問題なく出るはず。  

```bash
edinet search -c a会社,b会社
```

平均年収でも出せる。以下のクエリだと4000000以上の平均年収の会社を出す。

```bash
edinet search -s 4000000
```

まだORでつなげるだけなので、会社と平均年収をどちらもセットすると、どちらかに引っかかるものがすべて表示されてしまう。

## APIの構築

```bash
edinet api
```

localhost:8080に繋げればUIで表示される

# やりたいこと

- まだ従業員情報しか出していないので、収益などの情報も出せるようにしたい。
- 各引数ごとの検索はANDにしたい
- DBを消さずに更新できるようにしたい

# 考えメモ

一つのコマンドでいい気がしてきた。

- createdbコマンド -> metaのDBを作成
- searchコマンド -> 実際にsearchする。

処理の順番

createdbコマンド

- 書類のメタ情報を一定期間取ってくる
- それらのでーたのうち検索に必要なデータをDBに保存

その時の最新のデータのみをDBに入れておく。
いれる必要があるデータはまた考えておく

searchコマンド

- かいしゃ名をもらってそれを下にcreatecomdbがつくったDBを検索
- docIDからCSVデータをZIPで取得。
  - ファイルを展開して、情報を取得して、コマンドラインに表示。
  - (可能ならファイルを何処かに保存しておいて、その情報を取ってこれるようにすると処理は早くなりそう。)
