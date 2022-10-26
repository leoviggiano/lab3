package main

import (
	"fmt"
	"log"

	"blaus/pkg/entity"
	"blaus/pkg/http"
	"blaus/pkg/lib/csv"
)

const (
	QuantityToFetchRepositories = 100
)

func main() {
	requester, err := http.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	fetchPullRequests(requester)
}

func fetchPullRequests(requester http.Client) error {
	repositories, err := csv.OpenRepositories()
	if err != nil {
		return err
	}

	lastRepo := "elasticsearch"

	searchRepos := make([]*entity.Repository, 0)
	for i, r := range repositories {
		if r.Name == lastRepo {
			searchRepos = repositories[i:]
			break
		}
	}

	fmt.Println(searchRepos[0].Name)

	_, err = requester.QueryPullRequests(searchRepos)
	return err
}
