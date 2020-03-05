package scraper

import (
	"fmt"
	"time"

	"github.com/deshi-basara/kill-the-scout/database"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// Setup the colly collector used for scraping
func setup() *colly.Collector {
	// instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.immobilienscout24.de"),
		// colly.Debugger(&debug.LogDebugger{}),
		colly.MaxDepth(1),
	)

	// set random user agent on each request
	extensions.RandomUserAgent(c)

	return c
}

// ScrapeResultList performs the initial scraping of the result list and returns
func ScrapeResultList(url string) []string {
	c := setup()

	results := []string{}

	// on every element which has 'result-list'-class  call callback
	c.OnHTML(".result-list__listing", func(e *colly.HTMLElement) {
		exposeID := e.Attr("data-id")
		results = append(results, exposeID)
	})

	// set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// before making a request print
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping", r.URL.String())
	})

	// start scraping the handed url
	c.Visit(url)

	return results
}

// ScrapeExpose scrapes a single exposes data
func ScrapeExpose(exposeID string) database.Expose {
	c := setup()

	// build expose-url
	exposeURL := "https://www.immobilienscout24.de/expose/" + exposeID
	var expose database.Expose

	// on every element which has 'result-list'-class  call callback
	c.OnHTML("div[id=is24-content]", func(e *colly.HTMLElement) {
		title := e.ChildText("h1[data-qa]")
		address := e.ChildText(".is24-expose-address")
		rent := e.ChildText(".is24qa-kaltmiete.is24-value")
		size := e.ChildText(".is24qa-flaeche")
		rooms := e.ChildText(".is24qa-zi.is24-value")

		// build expose
		expose = database.Expose{
			ID:        exposeID,
			URL:       exposeURL,
			Title:     title,
			Address:   address,
			Rent:      rent,
			Size:      size,
			Rooms:     rooms,
			CreatedAt: time.Now(),
		}

		fmt.Printf("Expose scraped: %s\n", exposeID)
	})

	// before making a request print
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping", r.URL.String())
	})

	// start scraping the handed url
	c.Visit(exposeURL)

	return expose
}
