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
