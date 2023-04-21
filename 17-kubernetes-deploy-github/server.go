package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v52/github"
	"k8s.io/client-go/kubernetes"
)

type server struct {
	client           *kubernetes.Clientset
	githubClient     *github.Client
	webhookSecretKey string
}

func (s server) webhook(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	payload, err := github.ValidatePayload(req, []byte(s.webhookSecretKey))
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("unable to validate payload: %s", err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("unable to parse webhook: %s", err)
		return
	}
	switch event := event.(type) {
	case *github.Hook:
		fmt.Printf("found hook: %s\n", event)
	case *github.PushEvent:
		files := getFiles(event.Commits)
		fmt.Printf("found files: %s\n", strings.Join(files, ", "))
		for _, filename := range files {
			downloadFile, _, err := s.githubClient.Repositories.DownloadContents(ctx, *event.Repo.Owner.Name, *event.Repo.Name, filename, &github.RepositoryContentGetOptions{})
			if err != nil {
				w.WriteHeader(500)
				fmt.Printf("unable to download file: %s", err)
				return
			}
			defer downloadFile.Close()
			fileBody, err := io.ReadAll(downloadFile)
			if err != nil {
				w.WriteHeader(500)
				fmt.Printf("unable to read file: %s", err)
				return
			}
			_, _, err = deploy(ctx, s.client, fileBody)
			if err != nil {
				w.WriteHeader(500)
				fmt.Printf("unable to deploy: %s", err)
				return
			}
			fmt.Printf("deployed of %s successful\n", filename)
		}
	default:
		w.WriteHeader(500)
		fmt.Printf("unknown event type: %s", event)
		return
	}
}

func getFiles(commits []*github.HeadCommit) []string {
	allFiles := []string{}
	for _, commit := range commits {
		allFiles = append(allFiles, commit.Added...)
		allFiles = append(allFiles, commit.Modified...)
	}
	allUniqueFiles := make(map[string]bool)
	for _, filename := range allFiles {
		allUniqueFiles[filename] = true
	}
	allUniqueFilesSlice := []string{}
	for filename := range allUniqueFiles {
		allUniqueFilesSlice = append(allUniqueFilesSlice, filename)
	}
	return allUniqueFilesSlice
}
