package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"

	"server/config"
)

const (
	configpath = "config/config.toml"
)

var db = XormConn()

func GetDBConn() *xorm.Engine {
	return db
}

func XormConn() *xorm.Engine {

	db, err := xorm.NewEngine(GetDBConfig())
	if err != nil {
		fmt.Println("filed:connect db")
		panic(err.Error())
	}

	return db
}

func GetDBConfig() (string, string) {
	conf, err := config.ReadDBConfig(configpath)
	if err != nil {
		panic(err.Error())
	}

	CONNECT := conf.User + ":" + conf.Pass + "@" + conf.Protocol + "/" + conf.Dbname + "?parseTime=true&charset=utf8"

	return conf.Dbms, CONNECT
}
