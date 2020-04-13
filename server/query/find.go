package query

import (
	"server/model/join"
	"server/model/table"
	"time"

	"github.com/go-xorm/xorm"
)

func FindWorkInfos(db *xorm.Engine, date time.Time, userId uint64) ([]join.WorkInfos, error) {
	var workInfos []join.WorkInfos

	// wip MessyCode
	db.Where(
		"work_times.user_id = ?",
		userId,
	).And(
		"is_finished = true",
	).Iterate(&table.WorkTimes{}, func(idx int, bean interface{}) error {
		workTime := bean.(*table.WorkTimes)
		workInfos = append(workInfos, join.WorkInfos{
			WorkTimes: *workTime,
		})
		return nil
	})

	for i, w := range workInfos {
		db.Where(
			"work_time_id = ?",
			w.Id,
		).And(
			"is_finished = true",
		).Find(&workInfos[i].WorkRests)
	}

	return workInfos, nil
}

func FindWorkTimeLatest(db *xorm.Engine, cnt int, userId uint64) ([]table.WorkTimes, error) {
	var workTimes []table.WorkTimes

	err := db.Where(
		"user_id = ?",
		userId,
	).Limit(cnt).Desc(
		"started_at",
	).Find(workTimes)

	return workTimes, err
}
