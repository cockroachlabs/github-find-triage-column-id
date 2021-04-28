package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	token string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.token,
	}
	return token, nil
}

func getGithubClient() (*github.Client, error) {
	apiKey, ok := os.LookupEnv("GITHUB_API_KEY")
	if !ok {
		return nil, fmt.Errorf("cannot find GITHUB_API_KEY on local environment")
	}
	tokenSource := &tokenSource{
		token: apiKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	c := github.NewClient(oauthClient)
	return c, nil
}

var (
	flagGithubOrganization = flag.String(
		"organization",
		"cockroachdb",
		"organization containing the project",
	)
	flagGithubRepo = flag.String(
		"repo",
		"",
		"repository containing the project",
	)
	flagGithubProject = flag.String(
		"project",
		"",
		"project to lookup",
	)
	flagGithubColumn = flag.String(
		"column",
		"",
		"column to fetch from the project",
	)
)

func main() {
	flag.Parse()

	if *flagGithubOrganization == "" {
		log.Fatal("--organization must be specified")
	}

	if *flagGithubProject == "" {
		log.Fatal("--project must be specified")
	}

	if *flagGithubColumn == "" {
		log.Fatal("--column must be specified")
	}

	ghClient, err := getGithubClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	var r []*github.Project
	var resp *github.Response
	opts := &github.ProjectListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	more := true
	for more {
		if *flagGithubRepo != "" {
			r, resp, err = ghClient.Repositories.ListProjects(ctx, *flagGithubOrganization, *flagGithubRepo, opts)
		} else {
			r, resp, err = ghClient.Organizations.ListProjects(ctx, *flagGithubOrganization, opts)
		}
		if err != nil {
			log.Fatalf("error listing projects: %+v", err)
		}

		for _, proj := range r {
			if proj.GetName() == *flagGithubProject {
				findProjectColumn(ctx, ghClient, proj)
				return
			}
		}

		more = resp.NextPage != 0
		if more {
			opts.Page = resp.NextPage
		}
	}

	boardErrorSuffix := fmt.Sprintf("organization %s - maybe try specifying --repo", *flagGithubOrganization)
	if *flagGithubRepo != "" {
		boardErrorSuffix = fmt.Sprintf("repo %s/%s", *flagGithubOrganization, *flagGithubRepo)
	}
	log.Fatalf("unable to find project %s on %s", *flagGithubProject, boardErrorSuffix)
}

func findProjectColumn(ctx context.Context, ghClient *github.Client, proj *github.Project) {
	more := true
	opts := &github.ListOptions{
		PerPage: 100,
	}
	for more {
		cols, resp, err := ghClient.Projects.ListProjectColumns(ctx, proj.GetID(), opts)
		if err != nil {
			log.Fatalf("unable to list project columns: %+v", err)
		}
		for _, col := range cols {
			if col.GetName() == *flagGithubColumn {
				fmt.Printf("%d\n", col.GetID())
				return
			}
		}
		more = resp.NextPage != 0
		if more {
			opts.Page = resp.NextPage
		}
	}
	log.Fatalf("unable to find column %s on project %s", *flagGithubColumn, *flagGithubProject)
}
