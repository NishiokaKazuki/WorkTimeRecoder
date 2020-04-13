package query

import (
	"errors"
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func InsertUser(db *xorm.Engine, userName, hash string) (bool, error) {

	affected, err := db.Insert(&table.Users{
		Name: userName,
		Hash: hash,
	})

	return affected > 0, err
}

func InsertWorkTime(db *xorm.Engine, workTime table.WorkTimes) (bool, error) {

	affected, err := db.Insert(&workTime)
	return affected == 0, err
}

func InsertWorkRest(db *xorm.Engine, workRest table.WorkRests) (bool, error) {
	var r table.WorkRests

	has, _ := db.Where(
		"work_time_id = ?",
		workRest.WorkTimeId,
	).And(
		"is_finished = false",
	).Get(&r)
	if has != false {
		return false, errors.New("Should finish other resting.")
	}

	affected, err := db.Insert(&workRest)
	return affected == 0, err
}

func InsertSessionWorkTimes(db *xorm.Engine, sessionWorkTimes []table.SessionWorkTimes) (bool, error) {

	affected, err := db.Insert(&sessionWorkTimes)
	return affected == 0, err
}
