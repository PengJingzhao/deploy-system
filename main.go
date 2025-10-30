package main

import (
	"deploy-system/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许的源，可以是多个
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 是否允许发送 Cookie 等凭证
		MaxAge:           12 * time.Hour, // 预检请求缓存时间
	}))
	api.RegisterApi(engine)
	err := engine.Run("0.0.0.0:9099")
	if err != nil {
		log.Fatal(err)
	}
}
