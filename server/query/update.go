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
	).And(
		"disabled = false",
	).Update(&table.Users{
		Name: userName,
	})

	return affected > 0, err
}

func UpdateWorkTime(db *xorm.Engine, workTimes table.WorkTimes) (bool, error) {

	affected, err := db.Cols(
		"finished_at",
	).Where(
		"content = ?",
		workTimes.Content,
	).And(
		"user_id = ?",
		workTimes.UserId,
	).And(
		"disabled = false",
	).Update(&table.WorkTimes{
		FinishedAt: workTimes.FinishedAt,
	})

	return affected == 0, err
}

func UpdateWorkRest(db *xorm.Engine, workRest table.WorkRests) (bool, error) {

	affected, err := db.Cols(
		"is_finished",
		"finished_at",
	).Where(
		"work_time_id = ?",
		workRest.Id,
	).And(
		"disabled = false",
	).Update(&table.WorkRests{
		Isfinished: true,
	})

	return affected == 0, err
}
