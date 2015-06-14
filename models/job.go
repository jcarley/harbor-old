package models

import (
	"net/url"
	"os"
	"path/filepath"
)

type Job struct {
	LocalFolder  string
	RepoLocation string
}

func NewJob(localFolder string, repo string) Job {

	if localFolder == "" {
		localFolder = "Cellar"
	}

	repoUrl, _ := url.Parse(repo)

	root := localFolder
	host := repoUrl.Host
	account := filepath.Dir(repoUrl.Path)
	repoName := filepath.Base(repoUrl.Path)
	extension := filepath.Ext(repoName)
	name := repoName[0 : len(repoName)-len(extension)]

	wd, _ := os.Getwd()
	repoHome := filepath.Join(wd, root, host, account, name)

	return Job{repoHome, repo}
}
