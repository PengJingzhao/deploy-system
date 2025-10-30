package api

import (
	api "deploy-system/api/service"
	"github.com/gin-gonic/gin"
)

const prefix = "/app/v1"

func RegisterApi(group gin.IRouter) {
	router := group.Group(prefix)
	api.RegisterDeployService(router)
}
