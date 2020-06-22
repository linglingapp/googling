package googlescraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GoogleResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

// 選択可能なドメインの種類。
var googleDomains = map[string]string{
	"us": "https://www.google.com/search?hl=en&q=", // USA or default
	"jp": "https://www.google.co.jp/search?q=",     // Japan
	"uk": "https://www.google.co.uk/search?q=",     // United Kingdom
	"es": "https://www.google.es/search?q=",        // Spain
	"ca": "https://www.google.ca/search?q=",        // Canada
	"de": "https://www.google.de/search?q=",        // Deutschland
	"it": "https://www.google.it/search?q=",        // Italia
	"fr": "https://www.google.fr/search?q=",        // France
	"au": "https://www.google.com.au/search?q=",    // Australia
	"tw": "https://www.google.com.tw/search?q=",    // Taiwan
	"nl": "https://www.google.nl/search?q=",        // Nederland
	"br": "https://www.google.com.br/search?q=",    // Brasil
	"tr": "https://www.google.com.tr/search?q=",    // Turkey
	"be": "https://www.google.be/search?q=",        // Belgium
	"gr": "https://www.google.com.gr/search?q=",    // Greece
	"in": "https://www.google.co.in/search?q=",     // India
	"mx": "https://www.google.com.mx/search?q=",    // Mexico
	"dk": "https://www.google.dk/search?q=",        // Denmark
	"ar": "https://www.google.com.ar/search?q=",    // Argentina
	"ch": "https://www.google.ch/search?q=",        // Switzerland
	"cl": "https://www.google.cl/search?q=",        // Chile
	"at": "https://www.google.at/search?q=",        // Austria
	"kr": "https://www.google.co.kr/search?q=",     // Korea
	"ie": "https://www.google.ie/search?q=",        // Ireland
	"co": "https://www.google.com.co/search?q=",    // Colombia
	"pl": "https://www.google.pl/search?q=",        // Poland
	"pt": "https://www.google.pt/search?q=",        // Portugal
	"pk": "https://www.google.com.pk/search?q=",    // Pakistan
}

// 選択したドメインによるURL生成。選択していない場合，comにする。
func buildGoogleUrl(searchTerm string, countryCode string) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		return fmt.Sprintf("%s%s&num=100&hl=", googleBase, searchTerm)
	} else {
		return fmt.Sprintf("%s%s&num=100&hl=", googleDomains["us"], searchTerm)
	}
}

// 2ページ以降のURLを新たに取得。
func buildNextGoogleUrl(countryCode string, newUrl string) string {
	if googleBase, found := googleDomains[countryCode]; found {
		return fmt.Sprintf("%s%s", googleBase, newUrl)
	} else {
		return fmt.Sprintf("%s%s", googleDomains["us"], newUrl)
	}
}

func googleRequest(searchURL string) (*http.Response, error) {

	baseClient := &http.Client{}

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	time.Sleep(2 * time.Second) // sleep
	res, err := baseClient.Do(req)
	time.Sleep(2 * time.Second) // sleep
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// 2ページ目以降のURLのための新たな変数。
var newUrl string

func googleResultParser(response *http.Response) ([]GoogleResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []GoogleResult{}
	time.Sleep(2 * time.Second) // sleep

	// 次のページがあるのかどうかを確認。
	pageSel := doc.Find("table.AaVjTc")
	time.Sleep(2 * time.Second) // sleep
	for i := range pageSel.Nodes {
		pageNumber := pageSel.Eq(i)
		pageNumNow := pageNumber.Find("a#pnnext.G0iuSb")
		pageNumNow2, _ := pageNumNow.Attr("href")
		// 「次へ」のリンクがある場合，「newUrl」のアドレスを書き換える。
		if len(pageNumNow2) > 0 {
			pageNumNow3 := pageNumNow2[10:]
			newUrl = pageNumNow3
			time.Sleep(2 * time.Second) // sleep
		} else {
			newUrl = ""
			time.Sleep(2 * time.Second) // sleep
		}
	}

	// 「タイトル・URL・スニペット」を収集して「results」に格納。
	sel := doc.Find("div.g")
	rank := 1
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" && title != "" {
			result := GoogleResult{
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

func GoogleScrape(searchTerm string, countryCode string) ([]GoogleResult, error) {
	googleUrl := buildGoogleUrl(searchTerm, countryCode)
	res, err := googleRequest(googleUrl)
	time.Sleep(2 * time.Second) // sleep
	if err != nil {
		return nil, err
	}
	fmt.Println(res) // Checking '429 Too Many Requests'
	scrapes, err := googleResultParser(res)
	time.Sleep(2 * time.Second) // sleep

	var scrapes4 []GoogleResult
	var scrapes5 []GoogleResult
	fmt.Println(scrapes5)
	if newUrl != "" {
		// 2ページ目がある場合，新しいURLを生成。
		googleNewUrl := buildNextGoogleUrl(countryCode, newUrl)
		fmt.Println("Page 2:", googleNewUrl)
		res2, err := googleRequest(googleNewUrl)
		time.Sleep(2 * time.Second) // sleep
		if err != nil {
			return nil, err
		}

		scrapes2, err := googleResultParser(res2)
		scrapes4 = scrapes2
		time.Sleep(2 * time.Second) // sleep

		googleNewUrlLast := buildNextGoogleUrl(countryCode, newUrl)
		fmt.Println("Page 3:", googleNewUrlLast)
		res3, err := googleRequest(googleNewUrlLast)
		time.Sleep(2 * time.Second) // sleep
		if err != nil {
			return nil, err
		}
		scrapes3, err := googleResultParser(res3)
		scrapes5 = scrapes3
		time.Sleep(2 * time.Second) // sleep
	} else {
		goto EXIT
	}

EXIT:
	fmt.Println("csvファイルを作成しています…")

	// 結果をcsvファイルに書き込み，ダウンロード。
	file, err := os.Create("results.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"TITLE", "LINK", "SNIPPET"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, searchResult := range scrapes {
		searchResultSlice := []string{searchResult.ResultTitle, searchResult.ResultURL, searchResult.ResultDesc}
		srErr := w.Write(searchResultSlice)
		checkErr(srErr)
	}

	for _, searchResult2 := range scrapes4 {
		searchResultSlice2 := []string{searchResult2.ResultTitle, searchResult2.ResultURL, searchResult2.ResultDesc}
		srErr := w.Write(searchResultSlice2)
		checkErr(srErr)
	}

	for _, searchResult3 := range scrapes5 {
		searchResultSlice3 := []string{searchResult3.ResultTitle, searchResult3.ResultURL, searchResult3.ResultDesc}
		srErr := w.Write(searchResultSlice3)
		checkErr(srErr)
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
