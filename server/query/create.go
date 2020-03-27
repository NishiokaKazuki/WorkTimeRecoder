package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func CreateUser(db *xorm.Engine, userName, hash string) (bool, error) {

	affected, err := db.Insert(&table.Users{
		Name: userName,
		Hash: hash,
	})

	return affected > 0, err
}

func CreateWorkTime(db *xorm.Engine, workTimes table.WorkTimes) (bool, error) {

	affected, err := db.Insert(&workTimes)
	return affected > 0, err
}

func CreateWorkRest(db *xorm.Engine, workTimeId uint64) (bool, error) {

	affected, err := db.Insert(&table.WorkRests{
		WorkTimeId: workTimeId,
	})

	return affected > 0, err
}
