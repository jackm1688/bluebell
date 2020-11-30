package redis

import "fmt"

//redis
//redis尽量使用命名空间的方式区分不同的的Key,方便查询和拆分

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset;帖子及发帖的时间
	KeyPostScoreZSet   = "post:score"  //zset;帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset；记录用户及投票的类型 参数是post id
	KeyCommunitySetPF  = "community:"  //set;保存每个分区下帖子的id
)

// getRedisKey 给reids key加上前缘
func getRedisKey(key string) string {
	return fmt.Sprintf("%s%s", KeyPrefix, key)
}
