package controller

import (
	"Turing-Go/constant"
	"Turing-Go/server/common"
	"Turing-Go/server/web/logic"
	"Turing-Go/server/web/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var DefaultAccountController = &AccountController{}

type AccountController struct{}

func (a *AccountController) Register(ctx *gin.Context) {

	req := &model.RegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Println("解析参数异常:", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数错误"))
		return
	}
	err = logic.DefaultAccountLogic.Register(req)
	if err != nil {
		log.Println("注册业务异常:", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil))
}
