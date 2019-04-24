package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atlassian/go-artifactory/v2/artifactory"
	"github.com/atlassian/go-artifactory/v2/artifactory/transport"
	"github.com/atlassian/go-artifactory/v2/artifactory/v2"
	"log"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var username string
	var password string
	var url string
	var tokenName string
	var groups arrayFlags

	flag.StringVar(&username, "username", "admin", "username")
	flag.StringVar(&password, "password", "password", "password")
	flag.StringVar(&url, "url", "http://localhost:8080/artifactory", "artifactory URL")
	flag.StringVar(&tokenName, "name", "", "token username")
	flag.Var(&groups, "groups", "groups for the token")

	flag.Parse()

	if tokenName == "" {
		log.Fatalln("Must provide name for token")
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

	options := v2.AccessTokenOptions{
		Username: tokenName,
		Groups: groups,
		ExpiresIn: 0,
	}

	resp, err := rt.V2.Token.CreateAccessToken(context.Background(), &options)
	if err != nil {
		fmt.Printf("Error generating access token: %v\n", err)
		return
	}
	b, err := json.Marshal(resp)

	if err != nil {
		fmt.Printf("Error marshalling access token response: %v\n", err)
		return
	}

	fmt.Println("Response:")
	fmt.Println()
	fmt.Println(string(b))

	_, err = rt.V2.Token.RevokeAccessToken(context.Background(), resp.Token)
	if err != nil {
		fmt.Printf("Error revoking access token: %v\n", err)
		return
	}
}
