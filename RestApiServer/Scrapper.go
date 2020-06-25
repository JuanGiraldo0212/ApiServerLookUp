package main

import (
	_ "github.com/lib/pq"
	"log"
	"strings"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func scrapeResource(url string, item *Item) {
	// Request the HTML page.
	res, err := http.Get("https://www."+url)
	if err != nil {
	  log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
	  log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
	  log.Fatal(err)
	}
	title := doc.Find("title").Text()
	item.Title=title
	logo:=""
	head:=doc.Find("head link")
	head.EachWithBreak(func(index int, item *goquery.Selection) bool{
        linkTag := item
		link, _ := linkTag.Attr("rel")
        if strings.Contains(link,"icon"){
			logo,_=linkTag.Attr("href")
			return false
		}
		return true
	})
	item.Logo=logo
	
}


