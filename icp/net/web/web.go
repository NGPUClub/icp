package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nGPU/icp/configure"
	"github.com/nGPU/icp/middleware"
	"github.com/nGPU/icp/net/web/business"
	log4plus "github.com/nGPU/include/log4go"
)

type Web struct {
	webGin *gin.Engine
}

var gWeb *Web

func (w *Web) start() {
	funName := "start"
	log4plus.Info("start user gin listen")
	userGroup := w.webGin.Group("/user")
	{
		business.SingletonIcpERC20().Start(userGroup)
	}
	log4plus.Info("%s start Run Listen=[%s]", funName, configure.SingletonConfigure().Web.Listen)
	if err := w.webGin.Run(configure.SingletonConfigure().Web.Listen); err != nil {
		log4plus.Error("start Run Failed Not Use Http Error=[%s]", err.Error())
		return
	}
}

func SingletonWeb() *Web {
	if gWeb == nil {
		gWeb = &Web{}
		log4plus.Info("Create Web Manager")
		gWeb.webGin = gin.Default()
		gWeb.webGin.Use(middleware.Cors())
		gin.SetMode(gin.DebugMode)
		go gWeb.start()
	}
	return gWeb
}
