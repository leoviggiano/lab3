package queries

import (
	"encoding/json"
	"fmt"
)

type IssueStatus string

const (
	IssueOpen   = "OPEN"
	IssueClosed = "CLOSED"
)

type QueryOptions func(*graphqlRequest)
type graphqlRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func newRequest(query string) *graphqlRequest {
	return &graphqlRequest{
		Query:     query,
		Variables: make(map[string]string),
	}
}

func WithAfter(after string) func(*graphqlRequest) {
	return func(gql *graphqlRequest) {
		if len(after) == 0 {
			return
		}

		gql.Variables["after"] = after
	}
}

func WithIssueStatus(status IssueStatus) func(*graphqlRequest) {
	return func(gql *graphqlRequest) {
		gql.Variables["state"] = string(status)
	}
}

func Repositories(options ...QueryOptions) ([]byte, error) {
	query := `
	query search($after: String) {
		search(
			query: "is:public sort:stars-desc stars:>1000",
			type: REPOSITORY,
			first: 1,
      		after: $after
		) {
			pageInfo {
				endCursor
			}
			nodes {
				... on Repository {
					id
					name
					stargazerCount
					pullRequests {
						totalCount
					}
				}
			}
		}
	}`

	gqlRequest := newRequest(query)
	for _, option := range options {
		option(gqlRequest)
	}

	return json.Marshal(gqlRequest)
}

func PullRequests(repositoryID string, options ...QueryOptions) ([]byte, error) {
	query := fmt.Sprintf(`
	query search($after: String) {
		node(id: "%s") {
		  	... on Repository {
				pullRequests(
					first: 100,
					states: [MERGED, CLOSED],
					after: $after,
					orderBy: {field:CREATED_AT, direction: DESC}) {
						totalCount
						pageInfo {
							endCursor
						}
						nodes {
							id
							changedFiles
							createdAt
							updatedAt
							body
							additions
							deletions
							state
							closedAt
							assignees {
								totalCount
							}
							comments {
								totalCount
							}
							reviews(first: 10) {
								totalCount
								nodes {
									createdAt
									updatedAt
								}
							}
						}
				}
			}
		}
	}`, repositoryID)

	gqlRequest := newRequest(query)
	for _, option := range options {
		option(gqlRequest)
	}

	return json.Marshal(gqlRequest)
}
