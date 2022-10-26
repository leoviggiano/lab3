package entity

import (
	"fmt"
	"strconv"
	"time"
)

type Repository struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	StargazerCount int    `json:"stargazerCount"`
	After          string `json:"-"`

	NodePullRequests NodePullRequests `json:"pullRequests"`
}

func (r *Repository) CsvHeader() []string {
	return []string{"ID", "Name", "Stars", "Pull Requests", "After"}
}

func (r *Repository) CsvValues() []string {
	return []string{
		r.ID,
		r.Name,
		strconv.Itoa(r.StargazerCount),
		strconv.Itoa(r.NodePullRequests.TotalCount),
		r.After,
	}
}

func (r *Repository) FillFromCSV(row []string) {
	r.ID = row[0]
	r.Name = row[1]
	r.StargazerCount, _ = strconv.Atoi(row[2])
	r.NodePullRequests.TotalCount, _ = strconv.Atoi(row[3])
}

type PullRequest struct {
	ID             string        `json:"id"`
	State          string        `json:"state"`
	Body           string        `json:"body"`
	Additions      int           `json:"additions"`
	Deletions      int           `json:"deletions"`
	ChangedFiles   int           `json:"changedFiles"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	ClosedAt       time.Time     `json:"closedAt"`
	NodeReviews    NodeReviews   `json:"reviews"`
	NodeAssignees  NodeAssignees `json:"assignees"`
	NodeComments   NodeComments  `json:"comments"`
	RepositoryName string        `json:"-"`
}

func (p *PullRequest) CsvHeader() []string {
	return []string{"ID", "Repository", "State", "Body", "ChangedFiles", "CreatedAt", "UpdatedAt", "ClosedAt", "Reviews", "Assignees", "Comments", "Additions", "Deletions"}
}

func (p *PullRequest) CsvValues() []string {
	return []string{
		p.ID,
		p.RepositoryName,
		p.State,
		fmt.Sprint(len(p.Body)),
		fmt.Sprint(p.ChangedFiles),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		p.ClosedAt.String(),
		fmt.Sprint(p.NodeReviews.TotalCount),
		fmt.Sprint(p.NodeAssignees.TotalCount),
		fmt.Sprint(p.NodeComments.TotalCount),
		fmt.Sprint(p.Additions),
		fmt.Sprint(p.Deletions),
	}
}

type Review struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
