package main

import (
	
	"strings"
) 

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