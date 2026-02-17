package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	url := "https://books.toscrape.com/catalogue/see-america-a-celebration-of-our-national-parks-treasured-sites_732/index.html"

	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	// Extract main product information
	c.OnHTML("div.product_main", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		price := e.ChildText(".price_color")
		availability := strings.TrimSpace(e.ChildText(".availability"))

		fmt.Println("Title:", title)
		fmt.Println("Price:", price)
		fmt.Println("Availability:", availability)
	})		

	// Extract description
	c.OnHTML("#product_description", func(e *colly.HTMLElement) {
		description := e.DOM.Next().Text()
		fmt.Println("Description:", strings.TrimSpace(description))
	})

	// Extract table data
	c.OnHTML("table.table.table-striped tr", func(e *colly.HTMLElement) {
		key := e.ChildText("th")
		value := e.ChildText("td")

		fmt.Printf("%s: %s\n", key, value)
	})

	// Debug request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
