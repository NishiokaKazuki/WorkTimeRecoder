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

var db = xormConn()

func GetDBConn() *xorm.Engine {
	return db
}

func xormConn() *xorm.Engine {

	db, err := xorm.NewEngine(getDBConfig())
	if err != nil {
		fmt.Println("filed:connect db")
		panic(err.Error())
	}

	return db
}

func getDBConfig() (string, string) {
	conf, err := config.ReadDBConfig(configpath)
	if err != nil {
		panic(err.Error())
	}

	CONNECT := conf.User + ":" + conf.Pass + "@" + conf.Protocol + "/" + conf.Dbname + "?parseTime=true&charset=utf8"

	return conf.Dbms, CONNECT
}
