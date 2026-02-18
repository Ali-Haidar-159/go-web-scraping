package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Product struct {
	Name   string
	Price  string
	Rating string
	Image  string
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var products []Product

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://books.toscrape.com/"), // example site
		chromedp.WaitVisible("article.product_pod"),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll("article.product_pod")).map(p => {
				return {
					Name: p.querySelector("h3 a").getAttribute("title"),
					Price: p.querySelector(".product_price .price_color").innerText,
					Rating: p.querySelector(".star-rating").classList[1],
					Image: p.querySelector(".image_container img").src
				}
			})
		`, &products),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, p := range products {
		fmt.Println("Product", i+1)
		fmt.Println("Name:", p.Name)
		fmt.Println("Price:", p.Price)
		fmt.Println("Rating:", p.Rating)
		fmt.Println("Image:", p.Image)
		fmt.Println("-------------------")
	}
}
