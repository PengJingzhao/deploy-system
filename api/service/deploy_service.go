package api

import (
	"deploy-system/service"
	"github.com/gin-gonic/gin"
)

func RegisterDeployService(group gin.IRoutes) {
	group.POST("/deploy", service.DeployHandler)
}
