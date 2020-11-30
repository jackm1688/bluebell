package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//定义请求的参数结构体
type ParamSigup struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password"  binding:"required"`
	RePassword string `json:"confirm_password"  binding:"required,eqfield=Password"`
}

type ParamVoteData struct {
	//UserID //从当前请求中获取
	PostID    int64 `json:"post_id,string" binding:"required"`       //帖子ID
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` //赞成票(1)反对票(-1) 取消投票(0)
}

//获取帖子列表query string参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

type ParamCommunityPostList struct {
	*ParamPostList
	CommunityId int64 `json:"community_id" form:"community_id"`
}
