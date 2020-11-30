package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 用户投票处理函数
func PostVoteHandler(ctx *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseWithMsg(ctx, CodeInvalidParam, errData)
		return
	}
	//获取当用户的userID
	userID, err := getCurrentUserId(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	err = logic.VoteForPost(userID, p)
	if err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSucess(ctx, nil)
}
