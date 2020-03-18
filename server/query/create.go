package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func CreateWorkTime(db *xorm.Engine, userId uint64, content, supplement string) (bool, error) {

	affected, err := db.Insert(&table.WorkTimes{
		UserId:     userId,
		Content:    content,
		Supplement: supplement,
	})

	return affected > 0, err
}
