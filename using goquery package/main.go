package main 

import "fmt"
import "net/http"
import "github.com/PuerkitoBio/goquery"

const URL string = "https://books.toscrape.com/catalogue/soumission_998/index.html"

func main (){

	res,err := http.Get(URL)

	if(err != nil){
		fmt.Println("Find error to get the data : ",err)
		return 
	}

	defer res.Body.Close()

	if(res.StatusCode != 200){
		fmt.Println("The response status code is not 200. so the request is FAILED !")
		return 
	}

	doc,err := goquery.NewDocumentFromReader(res.Body)

	if(err != nil){
		fmt.Println("Find error to create the string HTML document : ",err) 
		return 
	}

	doc.Find("article table tbody tr").Each(func(i int, s *goquery.Selection){

		key := s.Find("th").Text()
		value := s.Find("td").Text()

		fmt.Printf("%s : %s : \n",key,value)
		
	})

}



