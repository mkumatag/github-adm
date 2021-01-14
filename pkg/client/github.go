package client

import (
	"context"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	PerPageSize = 100
)

type Github struct {
	*github.Client
	ctxt     context.Context
	Endpoint string
	APIKey   string
}

func NewGithub(baseURL, uploadURL, token string) (*Github, error) {
	var client *github.Client
	var err error
	var tc *http.Client

	ctxt := context.Background()

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctxt, ts)
	}

	if baseURL != "" && uploadURL != "" {
		client, err = github.NewEnterpriseClient(baseURL, uploadURL, tc)
		if err != nil {
			return nil, err
		}
	} else {
		client = github.NewClient(tc)
	}

	return &Github{
		ctxt:   ctxt,
		Client: client,
	}, nil
}

func (gh *Github) ListLabels(owner, repo string) ([]*github.Label, error) {
	var all []*github.Label
	opt := github.ListOptions{
		PerPage: PerPageSize,
	}
	for {
		labels, resp, err := gh.Issues.ListLabels(gh.ctxt, owner, repo, &opt)
		if err != nil {
			return nil, err
		}
		all = append(all, labels...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return all, nil
}

func (gh *Github) CreateLabel(owner, repo string, label *github.Label) (*github.Label, *github.Response, error) {
	return gh.Issues.CreateLabel(gh.ctxt, owner, repo, label)
}

func (gh *Github) DeleteLabel(owner, repo, name string) (*github.Response, error) {
	return gh.Issues.DeleteLabel(gh.ctxt, owner, repo, name)
}
