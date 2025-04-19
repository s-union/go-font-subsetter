package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"github.com/gocolly/colly"
)

// フィルタリング対象の不要な単語・フレーズ
var stopWords = []string{
	"メニュー", "検索", "ログイン", "コピーライト", "利用規約", "プライバシーポリシー",
}
var symbol = "、。，．・：；？！゛゜´｀¨＾￣＿ヽヾゝゞ〃仝々〆〇ー―‐／＼～∥｜…‥‘’“”（）〔〕［］｛｝〈〉《》「」『』【】＋－±×÷＝≠＜＞≦≧∞∴♂♀°′″℃￥＄￠￡％＃＆＊＠§☆★○●◎◇◆□■△▲▽▼※〒→←↑↓Å‰♯♭♪─│｡｢｣､①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩ㍉㌔㌢㍍㌘㌧㌃㌶㍑㍗㌍㌦㌣㌫㍊㌻㎜㎝㎞㎎㎏㏄㎡㍻〝〟№㏍℡㊤㊥㊦㊧㊨㈱㈲㈹㍾㍽㍼ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ０１２３４５６７８９abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()-^\\@[;:],./\\=~|`{+*}<>?_"
// 漢字、ひらがな、カタカナのみを抽出する関数
func filterJapaneseText(text string) string {
	re := regexp.MustCompile(`[一-龯ぁ-んァ-ンー]+`)
	matches := re.FindAllString(text, -1)
	return strings.Join(matches, " ")
}

// 文字列を正規化する関数（より強力な重複排除のため）
func cleanText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.Join(strings.Fields(text), " ")

	for _, word := range stopWords {
		text = strings.ReplaceAll(text, word, "")
	}

	// 漢字、ひらがな、カタカナのみを抽出
	text = filterJapaneseText(text)

	return text
}

func main() {
	startURL := "" // 文字を抽出したいWebサイトのURLを入力
	outputFile := "output.txt"

	c := colly.NewCollector(
		colly.AllowedDomains(""), // 許可するドメインを指定
		colly.MaxDepth(2), // 2階層までクロール
		colly.Async(true), // 非同期でクロール
	)

	visited := make(map[string]struct{}) // 訪問済みURLを管理
	charSet := make(map[rune]struct{})   // 使用されている文字を重複なく保存

	// ページ内のテキストを取得
	c.OnHTML("body", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)

		// 各文字を個別に処理して重複なく保存
		for _, char := range text {
			// 空白は無視
			if !unicode.IsSpace(char) {
				charSet[char] = struct{}{}
			}
		}
	})

	// リンクをたどってクロール範囲を拡大
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if _, found := visited[link]; !found {
			visited[link] = struct{}{}
			c.Visit(link)
		}
	})

	// クロール開始
	c.Visit(startURL)
	c.Wait()

	// 結果をファイルに保存
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("ファイルを作成できませんでした:", err)
		return
	}
	defer file.Close()

	// 文字をソートして保存
	chars := make([]rune, 0, len(charSet))
	for char := range charSet {
		chars = append(chars, char)
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	// 各文字を一行に出力
	for _, char := range chars {
		file.WriteString(string(char))
	}
	file.WriteString(string(symbol))
	fmt.Printf("サイトで使用されている文字を %s に保存しました（合計: %d 文字）\n", outputFile, len(charSet))
}
