package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func UpdateUser(db *xorm.Engine, userName, hash string) (bool, error) {

	affected, err := db.Cols(
		"name",
	).Where(
		"hash = ?",
		hash,
	).Update(&table.Users{
		Name: userName,
	})

	return affected > 0, err
}

func UpdateWorkTime(db *xorm.Engine, content string, userId uint64) (bool, error) {

	affected, err := db.Cols(
		"finished_at",
	).Where(
		"content = ?",
		content,
	).And(
		"user_id",
		userId,
	).Update(&table.WorkTimes{})

	return affected > 0, err
}

func UpdateWorkRest(db *xorm.Engine, workTimeId uint64) (bool, error) {

	affected, err := db.Cols(
		"finished_at",
	).Where(
		"work_time_id = ?",
		workTimeId,
	).Update(&table.WorkRests{})

	return affected > 0, err
}
