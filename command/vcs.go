package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jcarley/harbor/models"

	"golang.org/x/tools/go/vcs"
)

type VcsCommand struct {
}

func NewVcsCommand() *VcsCommand {
	return &VcsCommand{}
}

func (this *VcsCommand) Run(job models.Job) int {

	repo := job.RepoLocation
	localFolder := job.LocalFolder

	if _, err := os.Stat(localFolder); err != nil {
		log.Println(err)
	}

	if err := os.MkdirAll(localFolder, os.ModePerm); err != nil {
		log.Println(err)
	}

	privateFolder := filepath.Join(localFolder, ".git")

	cmd := vcs.ByCmd("git")

	if _, err := os.Stat(privateFolder); err != nil {
		fmt.Printf("Cloning %s ...\n", repo)
		cmd.Create(localFolder, repo)
		fmt.Printf("Finished cloning %s.\n", repo)
	} else {
		fmt.Printf("Updating %s ...\n", repo)
		cmd.Download(localFolder)
		fmt.Printf("Finished updating %s.\n", repo)
	}

	return 0
}
