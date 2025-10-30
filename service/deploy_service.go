package service

import (
	"deploy-system/deployer"
	"deploy-system/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeployHandler(ctx *gin.Context) {
	var req model.DeployReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	githubDeployer := deployer.NewGitHubDeployer(req.RepoURL, req.ImageName, req.ContainerName, req.Branch, req.PortMapping)
	err := githubDeployer.Deploy()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
