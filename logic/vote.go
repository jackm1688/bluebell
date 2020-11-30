package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"fmt"

	"go.uber.org/zap"
)

// 1.用户投票的数据
// 2.
// 投票功能
//投一票就加432分 86400/200 ->需要200张赞同可以给你的帖子续一天
/**
投票的几种情况:
direction=1时，两种情况:
 1.之前没投过票，现在投赞成票 -->更新分数和投票记录
 2.之前投反对票，现在投赞成票 -->更新分数和投票记录

direction=0时，两种情况:
 1.之前投过反对票，现在取消投票 -->更新分数和投票记录
 2.之前投过赞成票，现在取消投票 -->更新分数和投票记录
direction=-1,两种情况:
 1.之前没有投过票，现在投反对票   -->更新分数和投票记录
 2.之前投过赞成票，现在改投反对票 -->更新分数和投票记录

投票的限制:
每个帖子发表之日起一个星期内允许用户投票，超过用户就不允许用户投票了。
 1.到期之后将redis中保存的赞成票及对票存储到Mysql表中
 2.到期之后删除KeyPostVoteZSetPF
*/

// VoteForPost 为帖子投票函数
func VoteForPost(userId int64, p *models.ParamVoteData) error {
	zap.L().Debug("", zap.Any("userId", userId), zap.Any("postID", p.PostID))
	return redis.VoteForPost(fmt.Sprintf("%d", userId),
		fmt.Sprintf("%d", p.PostID),
		float64(p.Direction))
}
