package main

import "fmt"
import "time"

import "github.com/gocolly/colly"

import "ali.com/advanced_colly/config"
import "ali.com/advanced_colly/service"

func main(){

	collector := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
		colly.Async(true) ,
		colly.MaxDepth(2),
	)

	collector.Limit(&colly.LimitRule{
		DomainGlob : "*" ,
		Parallelism : 3 ,
		Delay : 1 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	collector.SetRequestTimeout(15 * time.Second)

	detailCollector := collector.Clone()

	// collector.SetProxy("http://127.0.0.1:8080")

	// all event handler 

	collector.OnRequest(func(req *colly.Request){

		ua := service.GetUserAgents()

		req.Headers.Set("User-Agent" , ua)

		// fmt.Println("Request is sending on , ",req.URL)
		// fmt.Println("User agent is : ",ua)
		// fmt.Println("")

	})

	collector.OnResponse(func(res *colly.Response){

		// fmt.Println("Status code : ",res.StatusCode,"\n\n")

	})


	// for extract data 

	collector.OnHTML("section article", func(element *colly.HTMLElement){

		title := element.ChildText("h3")

		fmt.Println(title)

	})
	
	collector.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
		// Product detail page link
		link := e.ChildAttr("div.image_container a", "href")
		fullLink := e.Request.AbsoluteURL(link)

		detailCollector.Visit(fullLink)
	})


	detailCollector.OnHTML("#product_description",func(e *colly.HTMLElement){

		data := e.ChildText("h2")
		fmt.Println(data)

	})

	// event listener for error 

	collector.OnError(func(res *colly.Response , err error){

		fmt.Println("Find error :",err)

	})

	// run or visit the collector and website 
	collector.Visit(config.URL)

	// wait until the work end 
	collector.Wait()
	detailCollector.Wait()

}