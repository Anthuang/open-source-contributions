package contributions

import "github.com/shurcooL/githubv4"

type PullRequestNode struct {
	Node struct {
		PullRequest struct {
			Title      string
			Url        string
			Repository struct {
				Name       string
				Url        string
				Visibility githubv4.RepositoryVisibility
			}
			Closed   bool
			ClosedAt githubv4.DateTime
		}
	}
}

type Viewer struct {
	ContributionsCollection struct {
		PullRequestContributions struct {
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
			Edges []PullRequestNode
		} `graphql:"pullRequestContributions(first: 100, after: $cursor)"`
	}
}
