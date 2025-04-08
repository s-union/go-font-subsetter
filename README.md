﻿# go-font-subsetter
 フォントのサブセット化のために、指定したWebサイトにおいて使用されている文字を抽出するGo製アプリケーション
## 使用方法
1. `main.go`の46行目の`startURL`に文字を抽出したいURLを入力
1. `main.go`の50行目の`colly.AllowedDomains("")`に許可するドメインを入力する
1. ターミナルで``` go run main.go ```
を実行する。するとoutput.txtに使用されている文字が抽出される。
### 注意
**ひらがな、カタカナ、漢字のみしか抽出されません！**
## 改良点
- [ ]　ローマ字と記号の抽出
  - output.txtに日常的に使用されるローマ字と記号を追加すればOKな気がする。
  - 関数filterJapaneseTextがそもそもいらない？
