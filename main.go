// Okteto CLI and preview environment functions

package main

import (
	"context"
	"encoding/json"
	"log"
	"okteto-dagger-module/internal/dagger"
	"strings"
)

type OktetoDaggerModule struct{}

// Define a struct to match the JSON structure for the endpoints object returned
type Endpoint struct {
	URL     string `json:"url"`
	Private bool   `json:"private"`
}

// Returns a container that has Okteto CLI with the correct context set

// example usage:
// dagger -m  call set-context --context=yourinstance.okteto.com --token=$OKTETO_TOKEN
func (m *OktetoDaggerModule) SetContext(context string, token string) *dagger.Container {
	return dag.Container().
		From("okteto/okteto").
		WithEnvVariable("OKTETO_TOKEN", token).
		// WithEnvVariable("OKTETO_CONTEXT", token).
		WithExec([]string{"okteto", "ctx", "use", context})
}

// Deploys a preview environment in the specified Okteto context

// example usage:
// dagger call preview-deploy --repo=https://github.com/RinkiyaKeDad/okteto-dagger-sample --branch=name-change --pr=https://github.com/RinkiyaKeDad/okteto-dagger-sample/pull/1 --context=yourinstance.okteto.com --token=$OKTETO_TOKEN
func (m *OktetoDaggerModule) PreviewDeploy(ctx context.Context,
	// Repo to deploy
	repo string,
	// Branch to deploy
	branch string,
	// URL of the pull request to attach in the Okteto Dashboard
	pr string,
	// Okteto context to be used for deployment
	context string,
	// Token to be used to authenticate with the Okteto context
	token string) (string, error) {
	c := m.SetContext(context, token).WithExec([]string{
		"okteto", "preview", "deploy", "--branch", branch, "--sourceUrl", pr, "--repository", repo, "--wait", strings.ToLower(branch),
	}).WithExec([]string{
		"okteto", "preview", "endpoints", strings.ToLower(branch), "--output=json",
	})

	endpointsOut, err := c.Stdout(ctx)
	if err != nil {
		return "", err
	}
	// Variable to hold the parsed data
	var endpoints []Endpoint

	// Parse the JSON data into the slice of Endpoint structs
	err = json.Unmarshal([]byte(endpointsOut), &endpoints)
	if err != nil {
		log.Fatal(err)
	}

	// StringBuilder to hold all URLs
	var urlsBuilder strings.Builder

	// Iterate through the parsed data and append each URL to the StringBuilder
	for _, endpoint := range endpoints {
		urlsBuilder.WriteString(endpoint.URL + "\n")
	}

	// Get the string with all URLs
	allURLs := urlsBuilder.String()

	return allURLs, nil
}

// Destroys a preview environment at the specified Okteto context

// example usage:
// dagger call preview-destroy --branch=name-change --context=yourinstance.okteto.com --token=$OKTETO_TOKEN
func (m *OktetoDaggerModule) PreviewDestroy(ctx context.Context,
	// Branch to deploy (to be used as the name for the preview env)
	branch string,
	// Okteto context to be used for deployment
	context string,
	// Token to be used to authenticate with the Okteto context
	token string) (string, error) {
	c := m.SetContext(context, token).WithExec([]string{
		"okteto", "preview", "destroy", strings.ToLower(branch), "--wait=false",
	})
	destoryOut, err := c.Stdout(ctx)
	if err != nil {
		return "", err
	}
	return destoryOut, nil
}
