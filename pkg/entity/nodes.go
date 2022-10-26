package entity

type PageInfo struct {
	EndCursor string `json:"endCursor"`
}

type SearchRepositories struct {
	Data struct {
		Search struct {
			PageInfo     PageInfo      `json:"pageInfo"`
			Repositories []*Repository `json:"nodes"`
		} `json:"search"`
	} `json:"data"`
}

type SearchPullRequest struct {
	Data struct {
		Node struct {
			NodePullRequests NodePullRequests `json:"pullRequests"`
		} `json:"node"`
	} `json:"data"`
}

type NodePullRequests struct {
	TotalCount   int           `json:"totalCount"`
	PageInfo     PageInfo      `json:"pageInfo"`
	PullRequests []PullRequest `json:"nodes"`
}

type NodeReviews struct {
	TotalCount int      `json:"totalCount"`
	Reviews    []Review `json:"nodes"`
}

type NodeAssignees struct {
	TotalCount int `json:"totalCount"`
}

type NodeComments struct {
	TotalCount int `json:"totalCount"`
}
