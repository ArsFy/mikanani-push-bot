package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var Config struct {
	Token  string `json:"token"`
	Admin  int64  `json:"admin"`
	Source string `json:"source"`
}

type ItemInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

var Data = struct {
	Tasks []string              `json:"tasks"`
	Items map[string][]ItemInfo `json:"items"`
}{
	Tasks: make([]string, 0),
	Items: make(map[string][]ItemInfo),
}

type ItemInfoLink struct {
	Url        string `json:"url"`
	Downloaded bool   `json:"downloaded"`
	Pushed     bool   `json:"pushed"`
}

var LinkMap = make(map[string]ItemInfoLink, 0)

func init() {
	// Config
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Failed to read config.json")
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalln("Failed to parse config.json")
	}

	// Data
	file, err = os.ReadFile("data.json")
	if err != nil {
		saveData()
	} else {
		err = json.Unmarshal(file, &Data)
		if err != nil {
			log.Fatalln("Failed to parse data.json")
		}
	}

	// Link
	file, err = os.ReadFile("link.json")
	if err != nil {
		saveLink()
	} else {
		err = json.Unmarshal(file, &LinkMap)
		if err != nil {
			log.Fatalln("Failed to parse data.json")
		}
	}
}

func saveData() {
	strJson, err := json.Marshal(Data)
	if err != nil {
		fmt.Println("Failed to marshal data.json")
		return
	}
	writeFile("data.json", strJson)
}

func saveLink() {
	strJson, err := json.Marshal(LinkMap)
	if err != nil {
		fmt.Println("Failed to marshal link.json")
		return
	}
	writeFile("link.json", strJson)
}
