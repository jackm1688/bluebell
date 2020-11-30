package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {

	dbCfg := settings.MySQLConfig{
		DriverName:   "mysql",
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "abc@198920",
		DBName:       "bluebell",
		MaxConns:     100,
		MaxIdleConns: 10,
	}

	st := &settings.AppConfig{}
	st.MySQL = &dbCfg
	_ = Init(st)
}
func TestCreatePost(t *testing.T) {

	post := &models.Post{
		ID:          10,
		AuthorID:    124,
		CommunityID: 1,
		Title:       "test",
		Content:     "this is a test",
	}
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("createPost failed:%v\n", err)
	} else {
		t.Logf("create sucessfully")
	}
}
