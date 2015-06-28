package main

import (
	"fmt"

	"github.com/jcarley/harbor/command"
	"github.com/jcarley/harbor/models"
)

func main() {
	server := NewServer()
	StartServer(server)
}

func testImageBuild() {

	ch := make(chan bool)

	repos := []string{"https://github.com/jcarley/docker-vault.git"}

	for _, repo := range repos {
		go buildImage(repo, ch)
	}

	select {
	case <-ch:
		fmt.Println("Finished!")
	}
}

func buildImage(repo string, finished chan bool) {

	workflow := []command.Command{
		command.NewVcsCommand(),
		command.NewBuildCommand(),
	}

	job := models.NewJob("", repo)
	for _, cmd := range workflow {
		cmd.Run(job)
	}
	finished <- true
}
