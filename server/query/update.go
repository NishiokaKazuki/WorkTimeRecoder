package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func UpdateWorkTime(db *xorm.Engine, content string) (bool, error) {

	affected, err := db.Cols(
		"finished_at",
	).Where(
		"content = ?",
		content,
	).Insert(&table.WorkTimes{})

	return affected > 0, err
}
