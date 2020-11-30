package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into post(post_id,author_id,community_id,title,content)values(?,?,?,?,?)"
	zap.L().Debug("", zap.String("SQL", sqlStr))
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return err
}

// GetPostById 根据帖子ID获取数据
func GetPostById(id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := "select post_id,author_id,community_id,status,title,content from post where post_id=?"
	zap.L().Debug("", zap.String("SQL", sqlStr), zap.Any("ID", id))
	err = db.Get(post, sqlStr, id)
	return
}

func GetPostList(offset, limit int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,status,title,content from 
         post ORDER BY create_time  DESC limit ?,?`
	zap.L().Debug("", zap.String("SQL", sqlStr))
	posts = make([]*models.Post, 0, limit)
	err = db.Select(&posts, sqlStr, offset, limit)
	return
}

// GetPostListByIds 根据指定的id列表查询帖子数据
func GetPostListByIds(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,status,title,content from 
         post  WHERE post_id IN (?) order by FIND_IN_SET(post_id,?)`
	zap.L().Debug("", zap.String("src_sql:", sqlStr))

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	zap.L().Debug("", zap.String("dst_sql:", query), zap.Any("args", args))
	db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
