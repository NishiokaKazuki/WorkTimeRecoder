package join

import "server/model/table"

type WorkInfos struct {
	table.WorkTimes `xorm:"extends"`
	WorkRests       []table.WorkRests `xorm:"extends"`
}
