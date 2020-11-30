package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ---- 跟社区相关的 ---

func CommunityHandler(ctx *gin.Context) {
	//1.查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic GetCommunityList failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSucess(ctx, data)
}

// 获取分类详情
func CommunityDeailHandler(ctx *gin.Context) {
	//1.获取社区id
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//2.根据社区的ID获取社区的详情
	communityDetail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("no community detail in db", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSucess(ctx, communityDetail)
}
