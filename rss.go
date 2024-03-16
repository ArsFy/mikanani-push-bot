package main

import (
	"encoding/xml"
	"log"
	"net/http"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Guid        string    `xml:"guid"`
	Link        string    `xml:"link"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Torrent     Torrent   `xml:"torrent"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type Torrent struct {
	Link          string `xml:"link"`
	ContentLength int    `xml:"contentLength"`
	PubDate       string `xml:"pubDate"`
}

type Enclosure struct {
	Type   string `xml:"type,attr"`
	Length int    `xml:"length,attr"`
	URL    string `xml:"url,attr"`
}

func getRss(uri string, fst bool) []ItemInfo {
	resp, err := http.Get("https://mikanani.me" + uri)
	if err != nil {
		log.Printf("Failed to fetch URL: %v", err)
		return nil
	}
	defer resp.Body.Close()

	var rss Rss
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		log.Fatalf("Failed to decode XML: %v", err)
		return nil
	}

	var downloaded = make([]ItemInfo, 0)
	if v, ok := Data.Items[rss.Channel.Title]; ok {
		downloaded = v
	}

	var need = make([]ItemInfo, 0)

	node := searchNode(downloaded)
	for _, item := range rss.Channel.Items {
		if !node.search(item.Title) {
			need = append(need, ItemInfo{
				Name: item.Title,
				Url:  item.Enclosure.URL,
			})
		}
	}

	if fst {
		Data.Tasks = append(Data.Tasks, uri)
		Data.Tasks = uniqueSlice(Data.Tasks)
	}

	downloaded = append(downloaded, need...)
	Data.Items[rss.Channel.Title] = downloaded
	saveData()

	return need
}
