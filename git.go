package main

import (
	"bytes"
	"os/exec"
	"regexp"
)

type GitClient struct {
	repo_path string
}

func (c *GitClient) Init(repo_path string) {
	c.repo_path = repo_path
}

func (c *GitClient) Status() ([]string, error) {
	cmd := exec.Command("git", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = c.repo_path
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	out_str := out.String()

	// TODO: Lift to Init
	// TODO: Find untracked files as well
	re, err := regexp.Compile(`(modified|new file):\s*(.*)`)
	if err != nil {
		return nil, err
	}

	results := re.FindAllStringSubmatch(out_str, -1)

	var matches = make([]string, len(results))
	for i, r := range results {
		matches[i] = r[2]
	}

	return matches, nil
}
