package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-rod/rod"
)

// Article represents the structure of our JSON output
type Article struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	PublishedOn string   `json:"published_on"`
	ArticleType string   `json:"article_type"`
	ArticleLink string   `json:"article_link"`
}

func scrape() {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Navigate to the target webpage
	page := browser.MustPage("https://www.frontiersin.org/search?query=bacteriophage&tab=articles&origin=https%3A%2F%2Fwww.frontiersin.org%2Fjournals").MustWaitLoad()
	log.Println("Browser started")

	// Extract the required data
	var articles []Article
	elements := page.MustElements("#article-results > ul > li")
	for _, el := range elements {
		title := el.MustElement("div.data-top > div.title").MustText()
		authorsEl := el.MustElements("div.data-top > ul > li")
		var authors []string
		for _, authorEl := range authorsEl {
			authors = append(authors, authorEl.MustText())
		}
		date := el.MustElement("div.data-bottom > div.date").MustText()
		articleType := el.MustElement("div.data-bottom > div.text > span.article-type").MustText()
		link := el.MustElement("a").MustProperty("href").String()

		articles = append(articles, Article{
			Title:       title,
			Authors:     authors,
			PublishedOn: date,
			ArticleType: articleType,
			ArticleLink: link,
		})
	}

	// Convert articles to JSON
	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", articles)
}

func main() {
	scrape()
}
