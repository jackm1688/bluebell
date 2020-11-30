package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Communtiy, err error) {
	sqlStr := "select community_id,community_name from community"
	if err = db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("this is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailById(id int64) (communityDeatail *models.CommunityDeatail, err error) {
	communityDeatail = new(models.CommunityDeatail)
	sqlStr := `select  community_id,community_name,
                introduction,create_time from community  where community_id=?`
	zap.L().Debug(sqlStr, zap.Any("ID", id))
	if err = db.Get(communityDeatail, sqlStr, id); err != nil {

		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDeatail, err

}
