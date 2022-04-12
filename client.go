package main

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func newClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}
