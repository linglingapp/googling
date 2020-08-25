package bingscraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BingResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

// 選択可能な地域コード
var bingCountryCode = map[string]string{
	"ar": "&cc=ar", //	Argentina
	"au": "&cc=au", //	Australia
	"at": "&cc=at", //	Austria
	"be": "&cc=be", //	Belgium
	"br": "&cc=br", //	Brazil
	"ca": "&cc=ca", //	Canada
	"cl": "&cc=cl", //	Chile
	"dk": "&cc=dk", //	Denmark
	"fi": "&cc=fi", //	Finland
	"fr": "&cc=fr", //	France
	"de": "&cc=de", //	Germany
	"hk": "&cc=hk", //	Hong Kong SAR
	"in": "&cc=in", //	India
	"id": "&cc=id", //	Indonesia
	"it": "&cc=it", //	Italy
	"jp": "&cc=jp", //	Japan
	"kr": "&cc=kr", //	Korea
	"my": "&cc=my", //	Malaysia
	"mx": "&cc=mx", //	Mexico
	"nl": "&cc=nl", //	Netherlands
	"nz": "&cc=nz", //	New Zealand
	"no": "&cc=no", //	Norway
	"cn": "&cc=cn", //	China
	"pl": "&cc=pl", //	Poland
	"pt": "&cc=pt", //	Portugal
	"ph": "&cc=ph", //	Philippines
	"ru": "&cc=ru", //	Russia
	"sa": "&cc=sa", //	Saudi Arabia
	"za": "&cc=za", //	South Africa
	"es": "&cc=es", //	Spain
	"se": "&cc=se", //	Sweden
	"ch": "&cc=ch", //	Switzerland
	"tw": "&cc=tw", //	Taiwan
	"tr": "&cc=tr", //	Turkey
	"gb": "&cc=gb", //	United Kingdom
	"us": "&cc=us", //	United States
}

var bingBaseUrl string = "https://www.bing.com/search?q="

// 選択した地域コードを使ってURL生成
func buildBingUrl(searchTerm string, countryCode string) string {
	// QUERYの前後にあるスペースを取り除く
	searchTerm = strings.Trim(searchTerm, " ")
	// QUERYに含まれている空白をすべて「+」に置き換える
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	// 入力した地域コードが目録にある場合
	if inputCountryCode, found := bingCountryCode[countryCode]; found {
		// パーセントエンコーディング（percent-encoding）によるエラーを回避
		searchTermEncoding := url.QueryEscape(searchTerm)
		// BingのURLにQUERYと地域コードを付ける
		fmt.Println("次のURLに接続（Bing）：", bingBaseUrl+searchTermEncoding+inputCountryCode)
		return fmt.Sprintf("%s%s%s", bingBaseUrl, searchTermEncoding, inputCountryCode)
	} else {
		// 地域コードを入力しなかった場合
		searchTermEncoding := url.QueryEscape(searchTerm)
		return fmt.Sprintf("%s%s", bingBaseUrl, searchTermEncoding)
	}
}

func bingRequest(searchURL string) (*http.Response, error) {
	// Clientを生成して用意したURLでリクエスト
	baseClient := &http.Client{}
	// ヘッダーを追加してGETリクエストを送信
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	time.Sleep(2 * time.Second) // sleep
	res, err := baseClient.Do(req)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// 2ページ目以降のURLのための新たな変数
var refinePartOfNextPageLink string
var newUrl string

func bingResultParser(response *http.Response) ([]BingResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	// BingResultタイプの配列生成
	results := []BingResult{}
	// 次のページがあるのかどうかを確認。
	findNextPage := doc.Find("li.b_pag")
	// ページにある内容をスクレイピング
	for i := range findNextPage.Nodes {
		nextPageButton := findNextPage.Eq(i)
		findNextPageLink := nextPageButton.Find("a.sb_pagN.sb_pagN_bp.sb_bp")
		partOfNextPageLink, _ := findNextPageLink.Attr("href")

		// 次ページへのリンクがある場合は取得したURLから必要のない部分を削除
		refinePartOfNextPageLink = partOfNextPageLink[10:]

		// Bingでページを遡る（前のページに移動する）のを防ぐ
		findNumberInNextPage := refinePartOfNextPageLink
		findNumberInPrePage := newUrl
		// 検索件数に関する数字を探索
		firstRegex := regexp.MustCompile(`=\d+&`)
		extractNumberAndSymbolNextP := firstRegex.Find([]byte(findNumberInNextPage))
		extractNumberAndSymbolPreP := firstRegex.Find([]byte(findNumberInPrePage))
		secondRegex := regexp.MustCompile(`\d+`)
		extractOnlyNumberNextP := secondRegex.Find([]byte(extractNumberAndSymbolNextP))
		extractOnlyNumberPreP := secondRegex.Find([]byte(extractNumberAndSymbolPreP))
		// 抽出した数字を整数に変換
		nextPstringToInt, _ := strconv.Atoi(string(extractOnlyNumberNextP))
		prePstringToInt, _ := strconv.Atoi(string(extractOnlyNumberPreP))
		// 次のページにある検索結果数が現在や前のページのそれより大きい場合
		if nextPstringToInt > prePstringToInt {
			// 次にスクレイピングするURLとして登録
			newUrl = refinePartOfNextPageLink
		} else {
			fmt.Println("これ以上進むことはできません（Bing）。")
			// 新しいURLがないので空けておく
			newUrl = ""
			break
		}
	}

	// 「タイトル・URL・スニペット」を収集して格納
	selectResult := doc.Find("li.b_algo")
	rank := 1
	for i := range selectResult.Nodes {
		item := selectResult.Eq(i)
		originalLink := item.Find("a")
		link, _ := originalLink.Attr("href")
		titleInResult := item.Find("h2")
		descInResult := item.Find("p")
		desc := descInResult.Text()
		title := titleInResult.Text()
		link = strings.Trim(link, " ")
		// タイトルやURLがない場合は収集しない
		if link != "" && link != "#" && title != "" {
			// 配列に追加
			result := BingResult{
				ResultRank:  rank,
				ResultURL:   link,
				ResultTitle: title,
				ResultDesc:  desc,
			}
			results = append(results, result)
			rank += 1
		}
	}
	return results, err
}

func BingScrape(searchTerm string, countryCode string) ([]BingResult, error) {
	// QUERYと地域コードでURLを生成
	bingUrl := buildBingUrl(searchTerm, countryCode)
	// サーバーにリクエスト
	res, err := bingRequest(bingUrl)
	time.Sleep(2 * time.Second) // sleep
	if err != nil {
		return nil, err
	}
	// fmt.Println(res) // '429 Too Many Requests' 確認用

	// スクレイピング開始
	scrapes, err := bingResultParser(res)
	time.Sleep(2 * time.Second) // sleep

	// csvファイルを作る
	file, err := os.Create("results.csv")
	fmt.Println("csvファイルを作成中")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()
	// ファイルのヘッダーを設定
	headers := []string{"TITLE", "URL", "SNIPPET"}

	wErr := w.Write(headers)
	checkErr(wErr)

	// スクレイピングした結果を書き込む
	for _, searchResult := range scrapes {
		searchResultSlice := []string{searchResult.ResultTitle, searchResult.ResultURL, searchResult.ResultDesc}
		srErr := w.Write(searchResultSlice)
		checkErr(srErr)
	}

	var newResults []BingResult
	// 新しいURLがある場合は上記と同じ操作を繰り返す
	for ok := true; ok; ok = (newUrl != "") {
		bingNewUrl := bingBaseUrl + newUrl
		fmt.Println("新しいURL（Bing）：", bingNewUrl)
		newRes, err := bingRequest(bingNewUrl)
		time.Sleep(2 * time.Second) // sleep
		if err != nil {
			return nil, err
		}
		newScrapes, err := bingResultParser(newRes)
		// 新たにスクレイピングしたものを配列に保存
		newResults = newScrapes
		if err != nil {
			return nil, err
		}
		// 保存したものをファイルに続けて書き込む
		for _, newSearchResult := range newResults {
			searchResultSlice2 := []string{newSearchResult.ResultTitle, newSearchResult.ResultURL, newSearchResult.ResultDesc}
			srErr := w.Write(searchResultSlice2)
			checkErr(srErr)
		}
	}
	if err != nil {
		return nil, err
	} else {
		return scrapes, nil
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
