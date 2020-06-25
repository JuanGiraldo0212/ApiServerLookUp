package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"log"
	"time"
	"strings"
)

func compareData(item *Item, uri string) {
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/server?sslmode=disable")
	if err != nil {
			log.Fatal("error connecting to the database: ", err)
	}

	items, err := db.Query("SELECT id, grade, time FROM items WHERE id = "+"'"+uri+"'")
    if err != nil {
        log.Fatal(err)
    }
    defer items.Close()
    for items.Next() {
				var id, grade string
				var t int64
        if err := items.Scan(&id, &grade, &t); err != nil {
            log.Fatal(err)
        }
				currTime:=time.Unix(t, 0)
				if (currTime.Add(time.Hour)).Before(time.Now()){
					item.Previous_ssl_grade=grade
					item.Servers_changed=false
					rows, err := db.Query("SELECT id, address, grade, country, owner FROM servers WHERE item = "+"'"+uri+"'")
    			if err != nil {
        		log.Fatal(err)
   				}
					defer rows.Close()
					i:=0
   				for rows.Next() {
        		var id, address, grade_s, country, owner string
        		if err := rows.Scan(&id, &address, &grade_s,&country,&owner); err != nil {
            log.Fatal(err)
						}
						currServer:=item.Endpoints[i]
						if currServer.IpAddress!=address || currServer.Grade!=grade_s || currServer.Country!=country || currServer.Owner!=owner{
							item.Servers_changed=true
						}
						i++
    			}

				} 
    }

}

func updateDB(item *Item, uri string){
	id:=uri
	grade:=item.Ssl_grade
	t:=time.Now().Unix()
	time := strconv.FormatInt(t, 10) 
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/server?sslmode=disable")
	if err != nil {
			log.Fatal("error connecting to the database: ", err)
	}
	if _, err := db.Exec(
		"UPSERT INTO items (id, grade, time) VALUES ("+"'"+id+"'"+", "+"'"+grade+"'"+", "+time+")"); err != nil {
		log.Fatal(err)
	}

	servers:=item.Endpoints
	for i:=0;i<len(servers);i++{
		index:=strconv.Itoa(i)
		curr:=servers[i]
		id_s:=uri+"_"+index
		address:=curr.IpAddress
		grade:=curr.Grade
		country:=curr.Country
		owner:=curr.Owner
		if _, err := db.Exec(
			"UPSERT INTO servers (id, address, grade, country, owner, item) VALUES ("+"'"+id_s+"'"+", "+"'"+address+"'"+", "+"'"+grade+"'"+", "+"'"+country+"'"+", "+"'"+owner+"'"+", "+"'"+id+"'"+")"); err != nil {
			log.Fatal(err)
		}
	}
}

func listDB() []string{
	var servers []string
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/server?sslmode=disable")
	if err != nil {
			log.Fatal("error connecting to the database: ", err)
	}

	rows, err := db.Query("SELECT id FROM items")
	if err != nil {
			log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
					log.Fatal(err)
			}
			servers=append(servers,(id+".com"))
	}
	return servers
}

func checkDB(item *Item,uri string){
	domain:=strings.Replace(uri, ".com", "", -1)
	compareData(item,domain)
	updateDB(item, domain)
}

