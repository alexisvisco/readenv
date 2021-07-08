package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type githubTag struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"commit"`
	NodeId string `json:"node_id"`
}

func checkVersion() {
	remoteVersion := fetchLatestTag()

	if remoteVersion == "" {
		return
	}

	currentVersion := getVersion()

	if currentVersion != remoteVersion {
		fmt.Printf("warning: readenv installed version is \u001b[31m%s\u001B[0m and the latest one is \u001b[32m%s\u001B[0m\n", currentVersion, remoteVersion)
		fmt.Printf("please consider upgrading by running \u001B[32mgo install github.com/alexisvisco/readenv@%s\n\u001b[0m", remoteVersion)
		fmt.Println("")
	}
}

func fetchLatestTag() (version string) {
	const fetchTagUrl = "https://api.github.com/repos/alexisvisco/readenv/tags?per_page=1"

	res, err := http.Get(fetchTagUrl)
	if err != nil {
		return ""
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	var tags []githubTag
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	if err := json.Unmarshal(body, &tags); err != nil {
		return ""
	}

	if len(tags) > 0 {
		return tags[0].Name
	}

	return ""
}
