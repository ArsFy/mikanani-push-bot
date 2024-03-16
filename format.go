package main

import (
	"strings"

	"golang.org/x/net/html"
)

func findClass(n *html.Node, class string) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode && containsClass(n, class) {
		nodes = append(nodes, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, findClass(c, class)...)
	}
	return nodes
}

func containsClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, class) {
			return true
		}
	}
	return false
}

func findLinks(n *html.Node) []string {
	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, findLinks(c)...)
	}
	return links
}

func findFirstAnchorText(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		return strings.TrimSpace(textContent(n))
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if text := findFirstAnchorText(c); text != "" {
			return text
		}
	}
	return ""
}

func textContent(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += textContent(c)
	}
	return text
}
