package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.AppConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host,
		cfg.MySQL.Port, cfg.MySQL.DBName)
	db, err = sqlx.Connect(cfg.MySQL.DriverName, dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(cfg.MySQL.MaxConns)
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)

	return
}

func Close() {
	_ = db.Close()
}
