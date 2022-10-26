package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"blaus/pkg/entity"
)

const (
	repositoriesPath = "etc/repositories.csv"
	pullRequestsPath = "etc/pull_requests.csv"
)

func Save(repositories map[string]*entity.Repository) error {
	if len(repositories) == 0 {
		return nil
	}

	file, err := os.OpenFile(repositoriesPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	values := make([][]string, 0, len(repositories))

	for _, r := range repositories {
		values = append(values, r.CsvValues())
	}

	if err := w.WriteAll(values); err != nil {
		return err
	}

	fmt.Printf("saved %d rows on csv with success\n", len(repositories))
	return nil
}

func SavePullRequests(pullRequests map[string]*entity.PullRequest) error {
	if len(pullRequests) == 0 {
		return nil
	}

	file, err := os.OpenFile(pullRequestsPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	values := make([][]string, 0, len(pullRequests))

	for _, r := range pullRequests {
		values = append(values, r.CsvValues())
	}

	if err := w.WriteAll(values); err != nil {
		return err
	}

	fmt.Printf("saved %d rows on csv with success\n", len(pullRequests))
	return nil
}

func OpenRepositories() ([]*entity.Repository, error) {
	file, err := os.Open(repositoriesPath)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	repositories := make([]*entity.Repository, 0)

	for _, row := range records[1:] {
		repository := &entity.Repository{}
		repository.FillFromCSV(row)

		repositories = append(repositories, repository)
	}

	return repositories, nil
}
