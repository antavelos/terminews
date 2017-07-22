package main

import (
	"github.com/antavelos/terminews/db"
	"log"
)

func main() {
	var data []db.RssReader
	var tdb *db.TDB
	var err error

	data = append(data, db.RssReader{Name: "Reuters World", Url: "http://feeds.reuters.com/Reuters/worldNews"})
	data = append(data, db.RssReader{Name: "Reuters Top News", Url: "http://feeds.reuters.com/reuters/topNews"})
	data = append(data, db.RssReader{Name: "Reuters Science", Url: "http://feeds.reuters.com/reuters/scienceNews"})
	data = append(data, db.RssReader{Name: "Reuters Technology", Url: "http://feeds.reuters.com/reuters/technologyNews"})
	data = append(data, db.RssReader{Name: "New York Times World", Url: "http://rss.nytimes.com/services/xml/rss/nyt/World.xml"})
	data = append(data, db.RssReader{Name: "New York Times Science", Url: "http://rss.nytimes.com/services/xml/rss/nyt/Science.xml"})
	data = append(data, db.RssReader{Name: "Eureca Alert, Breaking news", Url: "https://www.eurekalert.org/rss.xml"})

	if tdb, err = db.InitDB("term.db"); err != nil {
		log.Fatal(err)
	}
	if err = tdb.CreateTables(); err != nil {
		log.Fatal(err)
	}

	for _, rec := range data {
		if err = tdb.AddRssReader(rec); err != nil {
			log.Fatal(err)
		}
	}
}
