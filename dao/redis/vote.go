package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	oneWeekSeconds = 7 * 24 * 3600
	scorePerVote   = 432 //每一票占多少分
)

var (
	ErroVoteTimeExpire = errors.New("投票时间已经过")
	ErroVoteReplaced   = errors.New("不允许重复投票")
)

/**
投一票就加432分 86400/200 ->需要200张赞同可以给你的帖子续一天
投票的几种情况:
direction=1时，两种情况:
 1.之前没投过票，现在投赞成票 -->更新分数和投票记录 差值的绝对值: 1 +432
 2.之前投反对票，现在投赞成票 -->更新分数和投票记录 差值的绝对值: 2  +432*2

direction=0时，两种情况:
 1.之前投过反对票，现在取消投票 -->更新分数和投票记录 差值的绝对值: 1 -432
 2.之前投过赞成票，现在取消投票 -->更新分数和投票记录 差值的绝对值: 1 +432
direction=-1,两种情况:
 1.之前没有投过票，现在投反对票   -->更新分数和投票记录 差值的绝对值: 1 -432
 2.之前投过赞成票，现在改投反对票 -->更新分数和投票记录 差值的绝对值: 2 -432*2

投票的限制:
每个帖子发表之日起一个星期内允许用户投票，超过用户就不允许用户投票了。
 1.到期之后将redis中保存的赞成票及对票存储到Mysql表中
 2.到期之后删除KeyPostVoteZSetPF
*/

func VoteForPost(userId, postId string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postId).Val()
	if float64(time.Now().Unix())-postTime > oneWeekSeconds {
		return ErroVoteTimeExpire
	}

	//2.更新帖子的分数
	//先查当前帖子的投票记录
	zap.L().Debug("", zap.Any("key", getRedisKey(KeyPostVotedZSetPF+postId)),
		zap.Any("postId", postId), zap.Any("userId", userId))
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postId), userId).Val()
	if ov == value {
		return ErroVoteReplaced
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(ov - value) //计算两次投票的差值
	//更新分数
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postId)
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postId), postId)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postId), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userId,
		})
	}
	_, err := pipeline.Exec()
	return err
}

func CreatePost(postId int64, communityId int64) error {

	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityId)))
	pipeline.SAdd(cKey, postId)
	_, err := pipeline.Exec()
	return err
}
