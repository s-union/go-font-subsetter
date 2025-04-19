# go-font-subsetter
 フォントのサブセット化のために、指定したWebサイトにおいて使用されている文字を抽出するGo製アプリケーション
## 使用方法
1. `main.go`の46行目の`startURL`に文字を抽出したいURLを入力
1. `main.go`の50行目の`colly.AllowedDomains("")`に許可するドメインを入力する
1. ターミナルで``` go run main.go ```
を実行する。するとoutput.txtに使用されている文字が抽出される。

