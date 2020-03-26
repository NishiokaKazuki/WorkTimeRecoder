package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"

	"server/config"
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
	conf, err := config.ReadDBConfig()
	if err != nil {
		panic(err.Error())
	}

	CONNECT := conf.User + ":" + conf.Pass + "@" + conf.Protocol + "/" + conf.Dbname + "?parseTime=true&charset=utf8"

	return conf.Dbms, CONNECT
}
