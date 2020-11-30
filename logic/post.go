package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

// GetPostById 根据帖子ID获取数据
func GetPostById(id int64) (data *models.APIPostDeatail, err error) {
	//查询并组合数据
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById(id) failed",
			zap.Int64("id", id), zap.Error(err))
	}

	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed", zap.Int64("authorID", post.AuthorID),
			zap.Error(err))
		return
	}

	//根据社区id获取社区信息
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail failed",
			zap.Int64("communityId", post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.APIPostDeatail{
		Post:             post,
		CommunityDeatail: community,
	}
	data.AuthorName = user.Username
	return
}

// GetPostList 获取帖子列表
func GetPostList(offset, limit int64) (postList []*models.APIPostDeatail, err error) {
	posts, err := mysql.GetPostList(offset, limit)
	if err != nil {
		return nil, err
	}
	postList = make([]*models.APIPostDeatail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id获取社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed",
				zap.Int64("communityId", post.CommunityID),
				zap.Error(err))
			continue
		}
		data := &models.APIPostDeatail{
			Post:             post,
			CommunityDeatail: community,
		}
		data.AuthorName = user.Username
		postList = append(postList, data)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (postList []*models.APIPostDeatail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 3.根据id去数据库查询帖子详情
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	postList = make([]*models.APIPostDeatail, 0, len(posts))
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id获取社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed",
				zap.Int64("communityId", post.CommunityID),
				zap.Error(err))
			continue
		}
		data := &models.APIPostDeatail{
			Post:             post,
			CommunityDeatail: community,
			VoteNum:          voteData[idx],
		}
		data.AuthorName = user.Username
		postList = append(postList, data)
	}
	return
}

func GetCommunityPostList2(p *models.ParamCommunityPostList) (postList []*models.APIPostDeatail,
	err error) {

	// 2.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 3.根据id去数据库查询帖子详情
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	postList = make([]*models.APIPostDeatail, 0, len(posts))
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id获取社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed",
				zap.Int64("communityId", post.CommunityID),
				zap.Error(err))
			continue
		}
		data := &models.APIPostDeatail{
			Post:             post,
			CommunityDeatail: community,
			VoteNum:          voteData[idx],
		}
		data.AuthorName = user.Username
		postList = append(postList, data)
	}
	return
}
