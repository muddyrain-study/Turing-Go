package web

import (
	"Turing-Go/db"
	"Turing-Go/server/web/controller"
	"Turing-Go/server/web/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {

	db.TestDB()

	InitRouter(router)
}

func InitRouter(router *gin.Engine) {
	router.Use(middleware.Cors())
	router.Any("/account/register", controller.DefaultAccountController.Register)
}
