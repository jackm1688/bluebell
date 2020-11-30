package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Communtiy, err error) {
	//查找到所有community并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 根据ID查询分类详情
func GetCommunityDetail(id int64) (communityDeatail *models.CommunityDeatail, err error) {
	return mysql.GetCommunityDetailById(id)
}
