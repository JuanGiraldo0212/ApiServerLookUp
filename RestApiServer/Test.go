package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
)

type Server struct {
	IpAddress string `json:"address"`
	Grade string `json:"ssl_grade"`
	Country string `json:"country"`
	Owner string `json:"owner"`
}

type Item struct {
	Endpoints []*Server `json:"servers"`
	Servers_changed bool `json:"servers_changed"`
	Ssl_grade string `json:"ssl_grade"`
	Previous_ssl_grade string `json:"previous_ssl_grade"`
	Logo string `json:"logo"`
	Title string `json:"title"`
	Status bool `json:"is_down"`
}

func getServerInfo(uri string) []byte{
	body := doRequest("https://api.ssllabs.com/api/v3/analyze?host="+uri)
	item:=parseJson(string(body))
	servers:=item.Endpoints
	for i:=0;i<len(servers);i++{
		current:=servers[i]
		current.Country=string(doRequest("https://ipapi.co/"+item.Endpoints[0].IpAddress+"/country/"))
		current.Owner=string(doRequest("https://ipapi.co/"+item.Endpoints[0].IpAddress+"/org/"))
	}
	item.Ssl_grade=sslComparison(servers)
	scrapeResource(uri,item)
	item_marshalled, err := json.Marshal(item)
	if err==nil{
		return item_marshalled
	}
	return []byte("error")

}

func cleanJsonData(data string) string{
	newData:=strings.Split(data,":")
	curr:=strings.Trim(newData[1],`"`)
	return curr
}

func parseJson(json string) *Item{

trimData  := strings.Trim(json,"{}")
splitData:=strings.Split(trimData,`"endpoints":`)
fields := strings.Split(splitData[0],",")
status:=fields[4]
serversClean:=strings.Trim(splitData[1],"[]{}")
serversSplit:=strings.Split(serversClean,"},{")
servers:=[]*Server{}
for i:=0;i<len(serversSplit);i++{
	data:=strings.Split(serversSplit[i],",")
	var server *Server
	server=new(Server)
	server.IpAddress=cleanJsonData(data[0])
	server.Grade=cleanJsonData(data[3])
	servers=append(servers,server)
}

var item *Item
item=new (Item)
item.Endpoints=servers
item.Status=getStatus(cleanJsonData(status))

return item

}

func getStatus(status string) bool{
	if status=="READY"{
		return false
	}
	return true
}

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

func sslComparison(servers []*Server) string{
	grades:=make(map[string]int)
	grades["A+"]=7
	grades["A"]=6
	grades["B"]=5
	grades["C"]=4
	grades["D"]=3
	grades["E"]=2
	grades["F"]=1
	min:=99999
	grd:=""
	for i:=0;i<len(servers);i++{
		server:=servers[i]
		grade:=server.Grade
		point:=grades[grade]
		if point<min{
			min=point
			grd=grade
		}
	}
	return grd
}

func serverInfo(ctx *fasthttp.RequestCtx) {
	//fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	uri:=fmt.Sprintf("%v",ctx.UserValue("name"))
	fmt.Fprintln(ctx,string(getServerInfo(uri)))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/server/:name", serverInfo)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

func doRequest(url string) []byte {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)

	bodyBytes := resp.Body()
	//fmt.Printf(string(bodyBytes))
	return bodyBytes
	// User-Agent: fasthttp
	// Body:
}