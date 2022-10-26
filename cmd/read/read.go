package main

import (
	"fmt"
	"log"

	"blaus/pkg/lib/csv"
)

func main() {
	repositories, err := csv.OpenRepositories()
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range repositories {
		fmt.Printf("%#v\n", k.NodePullRequests)
	}

}
