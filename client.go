package main

import (
	"log"
	"time"

	"github.com/ryanlower/go-github/github" // Using fork because we need PR sorting
	"golang.org/x/oauth2"
)

type client struct {
	config *config
}

func (c *client) newGithubClient() *github.Client {
	token := &oauth2.Token{AccessToken: c.config.Github.AccessToken}
	client := new(oauth2.Config).Client(oauth2.NoContext, token)
	githubClient := github.NewClient(client)

	return githubClient
}

func (c *client) getPullRequestsMergedSince(since time.Time) []github.PullRequest {
	var mergedPRs []github.PullRequest
	var lastPRUpdated time.Time

	client := c.newGithubClient()

	opts := &github.PullRequestListOptions{
		State:       "closed",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 5},
	}
	for {
		PRs, resp, err := client.PullRequests.List(c.config.Repo.Owner,
			c.config.Repo.Name, opts)
		if err != nil {
			log.Fatal(err)
		}

		for _, PR := range PRs {
			if PR.MergedAt != nil && PR.MergedAt.After(since) {
				// Add to mergedPRs is this PR was merged after since
				mergedPRs = append(mergedPRs, PR)
			}

			lastPRUpdated = *PR.UpdatedAt
		}

		if resp.NextPage == 0 || lastPRUpdated.Before(since) {
			// Stop if there are no more pages or if the last PR we saw was updated
			// before since
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	return mergedPRs
}
