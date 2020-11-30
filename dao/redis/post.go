package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//3.ZRevRange 分数从大小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	/*for _,id := range ids{
		key := getRedisKey(KeyPostVotedZSetPF+id)
		//查询key中分数是1的元素数量->统计每篇帖子赞成票的数量
		v := client.ZCount(key,"1","1").Val()
		data = append(data,v)
	}*/

	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//3.ZRevRange 分数从大小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetCommunityPostIDsInOrder(p *models.ParamCommunityPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore把分区帖子set与帖子的分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据

	//社区的Key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityId)))
	//利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityId))
	if client.Exists(key).Val() < 1 {
		pipeline := client.Pipeline()
		//不存在，需要计算
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) //zinterstore计算
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
