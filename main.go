package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"
	"io/fs"
	"encoding/json"
)

var token string

//go:embed plan/*
var f embed.FS

type Gist struct {
	Description string						 `json:"description"`
	Public      bool						 `json:"public"`
	Files       map[string]map[string]string `json:"files"`
}

type Response struct {
	Url string `json:"url"`
}

func main() {
	isWeekly := false
	if len(os.Args) == 3 {
		isWeekly = true
	}
	files := readFiles(isWeekly)
	body := parseFiles(files, "План для студента", false)
	createGist(body)
}

func readFiles(isWeekly bool) []string {
    var fileNames []string

	files, err := fs.ReadDir(f, "plan")
	for _, file := range files {
        fileNames = append(fileNames, fmt.Sprintf("plan/%s", file.Name()))
    }
    if err != nil {
        panic(err)
    }

	if isWeekly {
		return fileNames[0:1]
	}
	return fileNames[1:]
}

func parseFiles(files []string, desc string, public bool) string {
	// parse files
	mFiles := make(map[string]map[string]string)
	filesLen := len(files)

    for _, file := range files {
		content, err := fs.ReadFile(f, file)
		if err != nil {
			log.Fatal(err)
		}

		splitFile := strings.Split(file, "/")
		filename := splitFile[len(splitFile)-1]

		contentStr := string(content)

		if len(os.Args) > 1 {
			desc = fmt.Sprintf("план для %s", os.Args[1])
		}
		
		if filesLen == 1 {
			if len(os.Args) > 2  {
				contentStr = strings.Replace(string(content), "<PLAN>", os.Args[2], 1)
			}
		}
 		
		mFiles[filename] = map[string]string{"content": contentStr}
    }

	gist := &Gist{Description: desc, Public: public, Files: mFiles}
	jsonGist, err := json.Marshal(gist)
    if err != nil {
        fmt.Printf("Error: %s", err.Error())
    } 

	return string(jsonGist)
}

func createGist(body string) {
	client := &http.Client{}
	var data = strings.NewReader(body)
	req, err := http.NewRequest("POST", "https://api.github.com/gists", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	if err := json.Unmarshal([]byte(bodyText), &response); err != nil {
        panic(err)
    }
	url := response.Url
	urlSplit := strings.Split(url, "/")
	id := urlSplit[len(urlSplit)-1]
	fmt.Printf("https://gist.github.com/Mariownyou/%s\n", id)
}
