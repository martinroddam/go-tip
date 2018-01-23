package gotip

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func getConfig() Config {
	var config Config
	filename := "./paths.yaml"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	ymlErr := yaml.Unmarshal(content, &config)
	if ymlErr != nil {
		log.Fatalf("error: %v", ymlErr)
	}

	return config
}
