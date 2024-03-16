package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
)

func uniqueSlice(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, str := range slice {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}

func writeFile(filename string, data []byte) {
	for {
		dst, err := os.Create(filename)
		if err != nil {
			namepath := strings.Split(filename, "/")
			os.MkdirAll(path.Join(strings.Join(namepath[:len(namepath)-1], "/")), os.ModePerm)
			continue
		}
		defer dst.Close()
		io.Copy(dst, bytes.NewReader(data))
		break
	}
}

type searchNodeType map[string]bool

func searchNode(list []ItemInfo) *searchNodeType {
	var nodeType = make(searchNodeType, len(list))
	for _, j := range list {
		nodeType[j.Name] = true
	}
	return &nodeType
}

func (it *searchNodeType) search(key string) bool {
	_, ok := (*it)[key]
	return ok
}
