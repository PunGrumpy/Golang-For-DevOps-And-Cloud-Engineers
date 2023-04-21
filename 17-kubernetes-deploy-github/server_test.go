package main

import (
	"testing"

	"github.com/google/go-github/v52/github"
)

func TestGetFiles(t *testing.T) {
	files := getFiles([]*github.HeadCommit{
		{
			ID:       github.String("123"),
			Added:    []string{"file1", "file2"},
			Modified: []string{},
			Message:  github.String("test commit"),
		},
	})
	if len(files) != len([]string{"file1", "file2"}) {
		t.Errorf("expected only 2 file: %+v", files)
	}
}
