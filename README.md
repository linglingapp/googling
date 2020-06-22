# Googling（更新中）

## 機能
Googleでの検索結果画面から，「タイトル」と「URL」が表示されている場合のみ，「タイトル・URL・スニペット」を一括ダウンロードします。「タイトル」のみ表示される場合もありますので，そのような例は除外します。現在は，最大300件をダウンロードするようになっています。

## インストール
プロジェクトディレクトリの中でビルドします（コードには「hello」となっている）。

```go
go build main.go
```

main.exeを実行するか，ターミナルなどで以下のコマンドを入力します。

```go
go run main.go  
```

## 使用方法
Chromeなどのウェブブラウザから「localhost:1323」に接続すると「クエリ（検索語句）」「地域」を入力するようになっています。クエリを入力した後，使用可能な地域コードを以下の中から選びます。何も入力しない場合，「com」で検索（英語で表示）します。日本語で検索する場合は「jp」を，韓国語で検索する場合は「kr」にした方がいいでしょう。

|     Country    | Code |
|:--------------:|:----:|
|       USA      |  us  |
|      Japan     |  jp  |
| United Kingdom |  uk  |
|      Spain     |  es  |
|     Canada     |  ca  |
|   Deutschland  |  de  |
|     Italia     |  it  |
|     France     |  fr  |
|    Australia   |  au  |
|     Taiwan     |  tw  |
|    Nederland   |  nl  |
|     Brasil     |  br  |
|     Turkey     |  tr  |
|     Belgium    |  be  |
|     Greece     |  gr  |
|      India     |  in  |
|     Mexico     |  mx  |
|     Denmark    |  dk  |
|    Argentina   |  ar  |
|   Switzerland  |  ch  |
|      Chile     |  cl  |
|     Austria    |  at  |
|      Korea     |  kr  |
|     Ireland    |  ie  |
|    Colombia    |  co  |
|     Poland     |  pl  |
|    Portugal    |  pt  |
|    Pakistan    |  pk  |

「検索＆ダウンロード」ボタンを押してしばらくすると，検索結果をcsvファイルでダウンロードします。ファイル名は「results」です。

## 注意点
APIは使っておりませんので，短い時間の間，頻繁に検索すると「429 Too Many Requests」を吐き出して検索できなくなります。むろん，そうなってもウェブブラウザでは普通に検索できます。429エラーが解除されるまでには，数時間かかります。

## 今後の計画
コードから重複を取り除いて，ウェブサーバーに配布する予定。
