package main

import (
	"testing"

	"github.com/google/go-github/v45/github"
)

func TestGetFiles(t *testing.T) {
	files := getFiles([]*github.HeadCommit{
		{
			Added:    []string{"test.txt"},
			Modified: []string{},
		},
	})
	if len(files) != 1 {
		t.Fatalf("Expected only 1 file: %+v", files)
	}
}
