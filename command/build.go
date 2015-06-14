package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jcarley/harbor/models"
)

type BuildCommand struct {
}

func NewBuildCommand() *BuildCommand {
	return &BuildCommand{}
}

func (this *BuildCommand) Run(job models.Job) int {

	path, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("installing fortune is in your future")
	}
	fmt.Printf("docker is available at %s\n", path)

	// docker build current directory
	cmdName := path
	location := job.LocalFolder
	cmdArgs := []string{"build", "--force-rm=true", "--no-cache=true", "-t", "jcarley/vault:latest", location}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("docker build out | %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	return 0
}
