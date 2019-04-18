package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/atlassian/go-artifactory/v2/artifactory"
	"github.com/atlassian/go-artifactory/v2/artifactory/transport"
	"io/ioutil"
	"os"
)

func main() {
	var username string
	var password string
	var url string
	var yamlFilename string

	flag.StringVar(&username, "username", "admin", "username")
	flag.StringVar(&password, "password", "password", "password")
	flag.StringVar(&url, "url", "http://localhost:8080/artifactory", "artifactory URL")
	flag.StringVar(&yamlFilename, "yaml", "", "yaml file to apply")

	flag.Parse()

	if len(yamlFilename) == 0 {
		fmt.Printf("Must specify YAML file to apply\n")
		os.Exit(1)
	}

	yaml, err := ioutil.ReadFile(yamlFilename)
	if err != nil {
		fmt.Printf("Error reading yaml file: %v\n", err)
		return
	}

	tp := transport.BasicAuth{
		Username: username,
		Password: password,
	}

	rt, err := artifactory.NewClient(url, tp.Client())
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	_, _, err = rt.V1.System.Ping(context.Background())
	if err != nil {
		fmt.Printf("Error testing authentication: %v\n", err)
		return
	} else {
		fmt.Println("Authenticated successfully")
	}

	_, err = rt.V2.Configuration.ApplyConfiguration(context.Background(), bytes.NewReader(yaml))
	if err != nil {
		fmt.Printf("Error applying configuration: %v\n", err)
		return
	} else {
		fmt.Println("Configuration applied successfully")
	}
}
