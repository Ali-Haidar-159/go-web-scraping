package main

import "fmt"
import "github.com/gocolly/colly/v2"

func main(){

	var URL string = "https://books.toscrape.com/catalogue/soumission_998/index.html"
	bookInfo := make(map[string]string)

	bookInfo["project"] = "Web Scraping"


	coll := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	// extract the heading 

	coll.OnHTML("div .product_main" , func (e *colly.HTMLElement){
		title := e.ChildText("h1")

		bookInfo["title"] = title

		// fmt.Println("The heading is : ",title)

	})


	// extract the table data and store in a map 
	coll.OnHTML("table tbody tr" , func (e *colly.HTMLElement){
		key := e.ChildText("th")
		data := e.ChildText("td")

		bookInfo[key] = data

		// fmt.Printf("%s : %s \n",key,data)
	})


	coll.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// Handle errors
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := coll.Visit(URL)
	if err != nil {
		fmt.Println(err)
	}


	defer func (){

		for key,data := range bookInfo{

			fmt.Println(key , " : " , data)

		}

	}()



}