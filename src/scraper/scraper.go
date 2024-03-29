package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

// Article represents the structure of our JSON output
type Article struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	PublishedOn string   `json:"published_on"`
	ArticleType string   `json:"article_type"`
	ArticleLink string   `json:"article_link"`
	Abstract    string   `json:"abstract"`
	Doi         string   `json:"doi"`
}

func extract_search_results(page *rod.Page, browser *rod.Browser) []Article {
	// Extract the required data
	var articles []Article
	elements := page.MustElements("#article-results > ul > li")
	for _, el := range elements {
		title := el.MustElement("div.data-top > div.title").MustText()
		// authorEl_Selector := fmt.Sprintf("#article-results > ul > li:nth-child(%d) > a > div.data-top > ul > li", i)
		authorsEl := el.MustElements("div.data-top > ul > li")
		var authors []string
		for _, authorEl := range authorsEl {
			authors = append(authors, authorEl.MustText())
		}
		date := el.MustElement("div.data-bottom > div.date").MustText()
		articleType := el.MustElement("div.data-bottom > div.text > span.article-type").MustText()
		// doiSelector := fmt.Sprintf("#article-results > ul > li:nth-child(%d) > a > div.data-bottom > ul > li.altmetric-embed > div", i)
		// doi := el.MustElement(doiSelector).MustText()
		// log.Printf("New page:")
		// log.Printf(browser.MustPages().Last().MustInfo().URL)

		// el.MustClick().Timeout(5 * time.Second).MustWaitStable()
		// newPage := browser.MustPage(browser.MustPages().Last().MustInfo().URL)
		// abstract := newPage.MustElement("#__layout > div > div.ArticlePage > div > div.Layout.Layout--withAside.Layout--withIbarMix.ArticleDetails > main > section > div.ArticleDetails__main__content > div > div.JournalFullText > div.JournalAbstract > p").MustText()
		// link := newPage.MustInfo().URL

		articles = append(articles, Article{
			Title:       title,
			Authors:     authors,
			PublishedOn: date,
			ArticleType: articleType,
			// ArticleLink: link,
			// Abstract:    abstract,
			// Doi:         doi,
		})
	}
	return articles
}

func scrape() {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()
	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	// Navigate to the target webpage
	// page := browser.MustPage("https://www.frontiersin.org/search?query=bacteriophage&tab=articles&origin=https%3A%2F%2Fwww.frontiersin.org%2Fjournals").MustWaitLoad()
	// log.Println("Browser started")

	page := browser.MustPage("https://www.frontiersin.org/search?tab=articles").MustWaitStable()

	page.MustElement("#search_query_input").MustInput("bacteriophage").MustType(input.Enter).MustWaitLoad()
	time.Sleep(5 * time.Second) // wait for 2 seconds

	articles := extract_search_results(page, browser)
	log.Printf("Articles: %s\n", articles)

	// Convert articles to JSON
	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		log.Fatal(err)
	}

	// Create a unique filename
	funcName := "scrape"
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.json", funcName, timestamp)

	// Open a new file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonArticles)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Data written to file: %s\n", filename)
}

// Shows how to disable headless mode and debug.
// Rod provides a lot of debug options, you can set them with setter methods or use environment variables.
// Doc for environment variables: https://pkg.go.dev/github.com/go-rod/rod/lib/defaults
func Example_disable_headless_to_debug() {
	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		Headless(false).
		Devtools(true)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	page := browser.MustPage("https://www.frontiersin.org/search?tab=articles")

	page.MustElement("#search_query_input").MustInput("bacteriophage").MustType(input.Enter)

	articles := extract_search_results(page, browser)

	fmt.Println(articles)

	utils.Pause() // pause goroutine
}

func main() {
	scrape()
}
