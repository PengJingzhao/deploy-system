package main

import (
	"deploy-system/deployer"
	"log"
)

func main() {
	repoURL := ""
	imageName := "hello"
	containerName := "hello-container"

	githubDeployer := deployer.NewGitHubDeployer(repoURL, imageName, containerName)
	githubDeployer.SetBranch("main")
	err := githubDeployer.Deploy()
	if err != nil {
		log.Fatal(err)
	}
}
