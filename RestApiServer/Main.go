package main

import (
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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

type Archive struct {
	Servers []string `json:"items"`
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
	item.Previous_ssl_grade=item.Ssl_grade
	scrapeResource(uri,item)
	checkDB(item,uri)
	item_marshalled, err := json.Marshal(item)
	if err==nil{
		return item_marshalled
	}
	return []byte("error")

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



func archive(ctx *fasthttp.RequestCtx){
	servers:=listDB()
	archive:=Archive{servers}
	data,_:=json.Marshal(archive)
	fmt.Fprintln(ctx,string(data))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/servers/:name", serverInfo)
	router.GET("/archive/", archive)
	log.Fatal(fasthttp.ListenAndServe(":8282", router.Handler))
}

func doRequest(url string) []byte {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)

	bodyBytes := resp.Body()
	return bodyBytes

}
