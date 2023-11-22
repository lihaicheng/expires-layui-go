package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/pkg/config"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/pkg/logger"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/server/controller"
	"go.uber.org/zap"
	"net/http"
)

func InitAPIServer(s *Server) error {
	s.Engine = gin.Default()
	zap.L().Info("setup router.")
	SetupRouter(s)
	err := Run(s)
	if err != nil {
		zap.L().Error("APIServer run failed", zap.Error(err))
	}
	return err
}

func SetupRouter(s *Server) {
	r := s.Engine
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("./layuimini/index.html")
	r.Static("/layuimini", "./layuimini")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/init", controller.DefaultInitMenuHandler)
	thingRoute := r.Group("/thing")
	thingRoute.POST("", CreateThing)
	thingRoute.GET("", ListThing)
	thingRoute.PUT("/:thing_id", UpdateThing)
	thingRoute.DELETE("/:thing_id", DeleteThing)

}

// Run engine
func Run(s *Server) error {
	addr := fmt.Sprintf("%s:%s", config.Config.APISetting.Host, config.Config.APISetting.Port)
	return s.Engine.Run(addr) // 对于HTTPS用RunTLS
}
