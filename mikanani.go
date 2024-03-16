package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func search(title string) (string, error) {
	resp, err := http.Get("https://mikanani.me/Home/Search?searchstr=" + url.QueryEscape(title))
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	linkList := make([]string, 0)

	nodes := findClass(doc, "list-inline an-ul")
	for _, node := range nodes {
		links := findLinks(node)
		linkList = append(linkList, links...)
	}

	if len(linkList) == 0 {
		return "", fmt.Errorf("not found")
	} else {
		return getAniRss(linkList[0])
	}
}

func getAniRss(uri string) (string, error) {
	resp, err := http.Get("https://mikanani.me" + uri)

	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	linkList := make([]string, 0)

	nodes := findClass(doc, "subgroup-text")
	for _, node := range nodes {
		title := findFirstAnchorText(node)
		if title == Config.Source {
			linkList = append(linkList, findLinks(findClass(node, "mikan-rss")[0])...)
		}
	}

	if len(linkList) == 0 {
		return "", fmt.Errorf("not found")
	} else {
		return linkList[0], nil
	}
}
