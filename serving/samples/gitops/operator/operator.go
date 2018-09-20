/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	ghclient "github.com/google/go-github/github"
	"github.com/knative/eventing/pkg/event"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/webhooks.v3/github"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// Environment variable containing json credentials
	envSecret = "GITHUB_SECRET"
)

// GithubHandler holds necessary objects for communicating with the Github.
type GithubHandler struct {
	client *ghclient.Client
	ctx    context.Context
}

type GithubSecrets struct {
	AccessToken string `json:"accessToken"`
	SecretToken string `json:"secretToken"`
}

func (h *GithubHandler) def(ctx context.Context, pl *github.PullRequestPayload) {
}

func (h *GithubHandler) newPullRequestPayload(ctx context.Context, pl *github.PullRequestPayload) {
	owner := pl.Repository.Owner.Login
	repo := pl.Repository.Name

	log.Printf("GOT PR with Title: %q", pl.PullRequest.Title)

	// Update the PR
	var pr *ghclient.PullRequest
	updatePr := func() (bool, error) {
		log.Printf("Updating PR")

		var err error
		pr, _, err = h.client.PullRequests.Get(ctx, owner, repo, int(pl.Number))

		if err != nil || pr.Mergeable == nil {
			return false, nil
		}

		return true, nil
	}

	err := wait.Poll(time.Second, 15*time.Second, updatePr)
	if err != nil {
		log.Fatalf("Could not fetch PR")
		return
	}

	log.Printf("Got updated PR: %q", *pr.Title)

	if *pr.State != "open" {
		log.Fatalf("PR not open")
		return
	}

	if !*pr.Mergeable {
		log.Fatalf("PR not mergeable")
		return
	}

	files, _, err := h.client.PullRequests.ListFiles(ctx, owner, repo, int(pl.Number), nil)
	if err != nil {
		log.Fatalf("Couldn't list files in PR: %q", err.Error())
		return
	}

	getFile := func(url string) (string, error) {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}

	for _, file := range files {
		log.Printf("Found file: %q", *file.Filename)

		log.Printf("Fetching content")
		content, err := getFile(*file.RawURL)
		if err != nil {
			log.Fatalf("Couldn't fetch file %v", err)
			return
		}

		log.Printf("Contents:\n%s", content)
	}

}

// look at PR
// can the PR be merged? if no, error
// apply config file
// was operation successful?
// if yes, merge pr
// if no, revert to old config

func main() {
	flag.Parse()
	githubSecrets := os.Getenv(envSecret)

	var credentials GithubSecrets
	err := json.Unmarshal([]byte(githubSecrets), &credentials)
	if err != nil {
		log.Fatalf("Failed to unmarshal credentials: %s", err)
		return
	}

	// Set up the auth for being able to talk to Github.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: credentials.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := ghclient.NewClient(tc)

	h := &GithubHandler{
		client: client,
		ctx:    ctx,
	}

	//log.Fatal(http.ListenAndServe(":8080", event.Handler(h.newPullRequestPayload)))
	log.Fatal(http.ListenAndServe(":8080", event.Handler(h.def)))
}
