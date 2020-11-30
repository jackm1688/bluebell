package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CratePostHandler(ctx *gin.Context) {
	//1.获取参数及参数校验
	p := new(models.Post)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("", zap.Any("ctx.ShouldBindJSON(p)", err))
		zap.L().Error("create post invalid param")
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//生成帖子Id
	p.ID = snowflake.GetID()
	authorId, err := getCurrentUserId(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	//获取帖子用户Id
	p.AuthorID = authorId
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSucess(ctx, nil)
}

func PostDetailHandler(ctx *gin.Context) {
	//获取参数
	postIdStr := ctx.Param("id")
	id, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(postIdStr,10,64) failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//根据帖子id取出数据
	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//返回响应数据
	ResponseSucess(ctx, data)
}

// GetPostListHandler 获取帖子列表处理函数
func GetPostListHandler(ctx *gin.Context) {
	//获取分页参数
	offset, limit := getPageInfo(ctx)
	//获取数据
	data, err := logic.GetPostList(offset, limit)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSucess(ctx, data)
}

// CratePostHandler 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(ctx *gin.Context) {
	//获取分页参数
	var p = &models.ParamPostList{
		Page:  0,
		Size:  10,
		Order: models.OrderTime,
	}
	err := ctx.ShouldBindQuery(p)
	if err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	//获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSucess(ctx, data)
}

//根据社区id获取帖子列表
func GetCommunityPostListHandler(ctx *gin.Context) {
	//获取分页参数
	var p = &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{
			Page:  0,
			Size:  10,
			Order: models.OrderTime,
		},
		CommunityId: 0,
	}
	err := ctx.ShouldBindQuery(p)
	if err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	//获取数据
	data, err := logic.GetCommunityPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostList2 failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSucess(ctx, data)
}

func getPageInfo(ctx *gin.Context) (offset, limit int64) {

	offsetStr := ctx.Query("page")
	limitStr := ctx.Query("size")
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	return offset, limit
}
