package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	containerName := ""
	if len(os.Args) == 1 {
		containerName = askContainerName()
	} else {
		containerName = os.Args[1]
	}
	containers := getContainers(containerName)
	options := []string{}
	containerNameByID := make(map[string]string)
	for _, container := range containers {
		options = append(options, container.Names[0])
		containerNameByID[container.Names[0]] = container.ID
	}
	selectedContainer := askContainer(options)
	shell := askShell()
	cmd := exec.Command("docker", "exec", "-ti", containerNameByID[selectedContainer], shell)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func askContainerName() string {
	name := ""
	err := survey.AskOne(
		&survey.Input{Message: "Search for container"},
		&name,
		survey.Required,
	)
	checkError(err)
	return name
}

func getContainers(name string) []types.Container {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	checkError(err)
	filters := filters.NewArgs()
	filters.Add("name", name)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters,
	})
	if len(containers) == 0 {
		log.Fatalln("No containers found.")
	}
	return containers
}

func askContainer(options []string) string {
	selectedContainer := ""
	err := survey.AskOne(
		&survey.Select{Message: "Select container", Options: options},
		&selectedContainer,
		survey.Required,
	)
	checkError(err)
	return selectedContainer
}

func askShell() string {
	shell := ""
	err := survey.AskOne(
		&survey.Input{Message: "Command", Default: "bash"},
		&shell,
		survey.Required,
	)
	checkError(err)
	return shell
}
