package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func DeleteSessionWorkTimes(db *xorm.Engine, userId uint64) error {
	var (
		s   table.SessionWorkTimes
		ids []uint64
	)

	err := db.Table("work_times").Cols("id").Where(
		"user_id = ?",
		userId,
	).Find(&ids)
	if err != nil {
		return err
	}

	_, err = db.In(
		"work_time_id = ?",
		ids,
	).Delete(&s)

	return err
}

func DeleteWorkTimes(db *xorm.Engine, workTimeId, userId uint64) (bool, error) {

	affected, err := db.Cols("disabled").Where(
		"Id = ?",
		workTimeId,
	).And(
		"user_id = ?",
		userId,
	).And(
		"disabled = false",
	).Update(&table.WorkTimes{
		Disabled: true,
	})

	return affected == 0, err
}
