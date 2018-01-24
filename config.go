package gotip

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func getPathInfo() PathInfo {
	var pathInfo PathInfo
	filename := "./paths.yaml"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	ymlErr := yaml.Unmarshal(content, &pathInfo)
	if ymlErr != nil {
		log.Fatalf("error: %v", ymlErr)
	}

	return pathInfo
}

func getGitInfo() GitInfo {
	var gitInfo GitInfo
	filename := "./secrets/secret.yaml"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	ymlErr := yaml.Unmarshal(content, &gitInfo)
	if ymlErr != nil {
		log.Fatalf("error: %v", ymlErr)
	}
	fmt.Println(gitInfo)
	return gitInfo
}
