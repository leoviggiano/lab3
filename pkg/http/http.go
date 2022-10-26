package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"blaus/pkg/config"
	"blaus/pkg/entity"
	"blaus/pkg/lib/csv"
	"blaus/pkg/queries"
)

type Client interface {
	QueryRepos(after string, limit int) (map[string]*entity.Repository, string, error)
	QueryPullRequests(repositories []*entity.Repository) (map[string]*entity.PullRequest, error)
}

type requester struct {
	client   *http.Client
	endpoint string
	token    string
}

func NewClient() (Client, error) {
	token := config.GithubToken()
	if len(token) == 0 {
		return nil, errors.New("empty github token")
	}

	return requester{
		client:   &http.Client{},
		endpoint: "https://api.github.com/graphql",
		token:    fmt.Sprintf("Bearer %s", token),
	}, nil
}

func (r requester) post(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", r.endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", r.token)
	req.Header.Add("Accept", "application/json")

	return r.client.Do(req)
}

func (r requester) QueryRepos(after string, limit int) (map[string]*entity.Repository, string, error) {
	repositories := make(map[string]*entity.Repository, limit)
	for len(repositories) < limit {
		query, err := queries.Repositories(queries.WithAfter(after))
		if err != nil {
			return nil, after, err
		}

		res, err := r.post(bytes.NewBuffer(query))
		if err != nil {
			return nil, after, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, after, err
		}

		parse := &entity.SearchRepositories{}
		err = json.Unmarshal(body, &parse)
		if err != nil {
			return nil, after, err
		}

		for _, v := range parse.Data.Search.Repositories {
			v.After = parse.Data.Search.PageInfo.EndCursor
			if _, ok := repositories[v.ID]; v.NodePullRequests.TotalCount >= 100 && !ok {
				repositories[v.ID] = v
				csv.Save(map[string]*entity.Repository{v.ID: v})
				fmt.Printf("Collected %d repositories\n", len(repositories))
			}
		}

		fmt.Println(after)
		after = parse.Data.Search.PageInfo.EndCursor
	}

	return repositories, after, nil
}

func (r requester) QueryPullRequests(repositories []*entity.Repository) (map[string]*entity.PullRequest, error) {
	pullRequests := make(map[string]*entity.PullRequest, 0)
	for _, repository := range repositories {
		after := ""

		prs := make(map[string]*entity.PullRequest)
		for len(prs) < 100 {
			query, err := queries.PullRequests(repository.ID, queries.WithAfter(after))
			if err != nil {
				return nil, err
			}

			res, err := r.post(bytes.NewBuffer(query))
			if err != nil {
				return nil, err
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			parse := &entity.SearchPullRequest{}
			err = json.Unmarshal(body, &parse)
			if err != nil {
				return nil, err
			}

			after = parse.Data.Node.NodePullRequests.PageInfo.EndCursor
			for _, p := range parse.Data.Node.NodePullRequests.PullRequests {
				fmt.Println(p.ID, repository.Name, p.NodeReviews.TotalCount)
				// if p.NodeReviews.TotalCount == 0 {
				// 	continue
				// }

				// firstReview := p.NodeReviews.Reviews[0]
				// if !firstReview.CreatedAt.After(p.CreatedAt.Add(time.Hour)) {
				// 	continue
				// }

				if _, ok := prs[p.ID]; !ok {
					p.RepositoryName = repository.Name
					prs[p.ID] = &p
					csv.SavePullRequests(map[string]*entity.PullRequest{p.ID: &p})
					fmt.Printf("Collected %d prs\n", len(prs))
				}
			}
		}

		for _, pr := range prs {
			pullRequests[pr.ID] = pr
		}
	}
	return pullRequests, nil
}

// func (r requester) queryIssueTotalCount(repositoryID string, status queries.IssueStatus) (int, error) {
// 	query, err := queries.Issues(repositoryID, queries.WithIssueStatus(status))
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := r.post(bytes.NewBuffer(query))
// 	if err != nil {
// 		return 0, err
// 	}

// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return 0, err
// 	}

// 	parse := &entity.NodeRepository{}
// 	err = json.Unmarshal(body, &parse)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return parse.Data.Node.Issues.TotalCount, nil
// }
