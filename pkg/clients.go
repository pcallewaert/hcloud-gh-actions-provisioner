package pkg

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/oauth2"
)

func setupHcloudClient(hcloudToken string) *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(hcloudToken))
}

func setupGithubClient(githubPat string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubPat},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}
