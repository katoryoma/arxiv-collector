package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// ■ 1. 型定義 (PythonのDataclassに近い)
// arXivのAPIはXML形式で返ってくるので、その構造を事前に定義します。
// Pythonと違い、Goは「どんなデータが来るか」を厳密に決める必要があります。

type Feed struct {
	Entry []Entry `xml:"entry"` // <entry>タグの中身をリストで持つ
}

type Entry struct {
	Title      string   `xml:"title"`                                     // <title>タグ
	ID         string   `xml:"id"`                                        // <id>タグ (URL)
	Summary    string   `xml:"summary"`                                   // <summary>タグ
	Published  string   `xml:"published"`                                 // <published>タグ (発表日)
	Authors    []Author `xml:"author"`                                    // <author>タグのリスト
	JournalRef string   `xml:"http://arxiv.org/schemas/atom journal-ref"` // ジャーナル情報
}

type Author struct {
	Name string `xml:"name"` // 著者名
}

// ■ 2. 論文データの入れ物
type Paper struct {
	Keyword     string
	Title       string
	URL         string
	Summary     string
	Authors     string
	Published   string
	Publication string // 出版情報（ジャーナル、会議、等）
}

func main() {
	for {
		start := time.Now()

		// ユーザーからキーワード入力を受け取る（複数単語対応）
		fmt.Println("検索したいキーワードを入力してね (複数単語可):")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input) // 改行とスペースをトリム

		// 論文数の入力
		fmt.Println("いくつみつけてほしいの？100までにしてね！:")
		maxResultsStr, _ := reader.ReadString('\n')
		maxResultsStr = strings.TrimSpace(maxResultsStr)
		maxResults := 100
		if _, err := fmt.Sscanf(maxResultsStr, "%d", &maxResults); err != nil || maxResults <= 0 || maxResults > 100 {
			maxResults = 100
		}

		// 年数の入力（オプション）は廃止 - arXIvの仕様上、正確な年度フィルタリングが難しいため
		fmt.Println("⚠️  年度フィルタリングは非対応です。最新の論文を取得します")

		// 入力されたキーワード（複数の場合はカンマ区切り）
		keywords := []string{input}

		// 結果を受け取るための「チャネル」
		// PythonのQueueのようなもので、並列処理中のデータの通り道です。
		results := make(chan []Paper, len(keywords))

		// WaitGroup: 並列処理の管理係。「全員終わるまで待つよ」というカウンター。
		var wg sync.WaitGroup

		fmt.Println("🚀 たくさん収集するね！")

		// ■ 3. 並行処理の開始 (ここがGoの真骨頂！)
		for _, kw := range keywords {
			wg.Add(1) // 「仕事が1つ増えた」とカウント

			// "go func()..." と書くだけで、別のスレッド(Goroutine)で同時に走ります。
			go func(k string) {
				defer wg.Done() // 終わったら「仕事終わった」と報告
				data := fetchPapers(k, maxResults)
				results <- data // 結果をチャネル（通り道）に投げる
			}(kw)
		}

		// 全員の仕事が終わるのを待って、チャネルを閉じるための監視係
		go func() {
			wg.Wait()
			close(results)
		}()

		// ■ 4. CSVファイルへの書き出し
		// ファイル名を「日時_キーワード.csv」の形式で生成
		timestamp := time.Now().Format("20060102_150405")
		// ファイル名用にスペースをアンダースコアに置き換え
		filenameKeyword := strings.ReplaceAll(input, " ", "_")
		filename := fmt.Sprintf("%s_%s.csv", timestamp, filenameKeyword)
		file, _ := os.Create(filename)
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// ヘッダー書き込み
		writer.Write([]string{"Keyword", "Title", "URL", "Summary", "Authors", "Published", "Publication"})

		// チャネルから結果が流れてくるたびに書き込む
		count := 0
		for papers := range results {
			for _, p := range papers {
				writer.Write([]string{p.Keyword, p.Title, p.URL, p.Summary, p.Authors, p.Published, p.Publication})
				count++
			}
		}

		fmt.Printf("✅ %d 件の論文を保存しといたわ (%s)///\n", count, filename)
		fmt.Printf("⏱ 頑張った時間: %v\n", time.Since(start))

		// 続けるかどうか確認
		fmt.Println("\n他にも探したい？(y/n):")
		reader = bufio.NewReader(os.Stdin)
		continueStr, _ := reader.ReadString('\n')
		continueStr = strings.TrimSpace(continueStr)
		if continueStr != "y" && continueStr != "Y" {
			fmt.Println("✨ すぐきてね？	バイバイ♡")
			break
		}
	}
}

// arXiv APIを叩いてデータを取得する関数
func fetchPapers(keyword string, maxResults int) []Paper {
	// URLエンコード (スペースを %20 にするなど)
	q := url.QueryEscape("all:" + keyword)
	apiURL := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=%s&start=0&max_results=%d&sortBy=submittedDate&sortOrder=descending", q, maxResults)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	// XMLデータを読み込む
	body, _ := io.ReadAll(resp.Body)

	// 事前に定義したStruct(型)にデータを流し込む (Unmarshall)
	var feed Feed
	xml.Unmarshal(body, &feed)

	var papers []Paper
	yearCounts := make(map[string]int) // 年度別の件数をカウント

	for _, entry := range feed.Entry {
		// 論文の年度を取得
		yearStr := ""
		if len(entry.Published) >= 4 {
			yearStr = entry.Published[:4]
			yearCounts[yearStr]++
		}

		// 著者名をセミコロン區切りで結合
		authors := ""
		for i, author := range entry.Authors {
			if i > 0 {
				authors += "; "
			}
			authors += author.Name
		}

		papers = append(papers, Paper{
			Keyword:     keyword,
			Title:       entry.Title,
			URL:         entry.ID,
			Summary:     entry.Summary,
			Authors:     authors,
			Published:   entry.Published,
			Publication: entry.JournalRef,
		})
	}
	return papers
}
