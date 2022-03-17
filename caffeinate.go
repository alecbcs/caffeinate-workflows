package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/alecbcs/caffeinate-workflows/config"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func main() {
	// Print out the little caffinate icon in a terrible way to
	// prevent gofmt from converting spaces into tabs.
	fmt.Printf(
"           ((((\n" +
"          ((((\n" +
"           ))))\n" +
"        _ .---.\n" +
"       ( |'---'|\n" +
`        \|     |` + "\n" +
"        : .___, :\n" +
"         '-----'\n" +
"Caffeinate-Workflow %s \n\n", config.Version)

	// Parse environment variables into config.
	config, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize context for Git client
	ctx := context.Background()

	// Generate oauth token/client from config.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHub.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Initialize GitHub client.
	client := github.NewClient(tc)

	// Get GitHub repository owner and name from path string.
	pathData := strings.Split(config.GitHub.Repository, "/")
	owner := pathData[0]
	repo := pathData[1]

	var files []string

	if len(config.Workflow.Files) > 0 {
		// Extract workflow files names from config
		files = strings.Split(config.Workflow.Files, ",")
	} else {
		// If no workflow files were defined assume action should apply to
		// all workflows within the repository.
		workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repo, nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, workflow := range workflows.Workflows {
			files = append(files,
				filepath.Base(workflow.GetPath()),
			)
		}
	}

	// Enable each workflow file in config.
	for _, file := range files {
		_, err := client.Actions.EnableWorkflowByFileName(
			ctx,
			owner,
			repo,
			file,

		)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Enabled: %s\n", file)
	}

	fmt.Println("Done!")
}
