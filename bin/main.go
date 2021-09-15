package main

import (
	"context"
	fmt "fmt"
	"os"

	"github.com/anthuang/open-source-contributions/internal/contributions"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type RepoToPRs map[string][]contributions.PullRequestNode

func publicPRs(prs []contributions.PullRequestNode) []contributions.PullRequestNode {
	var publicPrs []contributions.PullRequestNode
	for _, pr := range prs {
		if pr.Node.PullRequest.Repository.Visibility == githubv4.RepositoryVisibilityPublic {
			publicPrs = append(publicPrs, pr)
		}
	}
	return publicPrs
}

func printOutput(repoToPRs RepoToPRs) {
	fmt.Println("## Open Source Contributions")
	for repo, prs := range repoToPRs {
		fmt.Printf("%s\n", repo)

		for _, pr := range prs {
			fmt.Printf(
				"- [%s](%s)\n",
				pr.Node.PullRequest.Title,
				pr.Node.PullRequest.Url,
			)
		}
		fmt.Println("")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"cursor": (*githubv4.String)(nil),
	}
	repoToPRs := RepoToPRs{}

	var query struct {
		Viewer contributions.Viewer
	}

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		publicPRs := publicPRs(query.Viewer.ContributionsCollection.PullRequestContributions.Edges)
		for _, pr := range publicPRs {
			if !pr.Node.PullRequest.Closed {
				continue
			}
			repoToPRs[pr.Node.PullRequest.Repository.Name] = append(repoToPRs[pr.Node.PullRequest.Url], pr)
		}

		pageInfo := query.Viewer.ContributionsCollection.PullRequestContributions.PageInfo
		if !pageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(pageInfo.EndCursor)
	}

	printOutput(repoToPRs)
}
