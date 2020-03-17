package table

import (
	"time"
)

type Users struct {
	Id uint64
	Name string
}

type WorkTimes struct {
	Id uint64
	UserId uint64
	Content string
	Supplement string
	StartedAt time.Time
	FinishedAt time.Time
}

type WorkRests struct {
	Id uint64
	WorkTimeId uint64
	StartedAt time.Time
	FinishedAt time.Time
}