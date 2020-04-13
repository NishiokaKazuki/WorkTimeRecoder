package table

import (
	"time"
)

type Users struct {
	Id   uint64
	Name string
	Hash string
}

type WorkTimes struct {
	Id         uint64
	UserId     uint64
	Content    string
	Supplement string
	Isfinished bool
	StartedAt  time.Time
	FinishedAt time.Time
}

type WorkRests struct {
	Id         uint64
	WorkTimeId uint64
	Isfinished bool
	StartedAt  time.Time
	FinishedAt time.Time
}

type SessionWorkTimes struct {
	WorkTimeId uint64
	Hash       string
}
