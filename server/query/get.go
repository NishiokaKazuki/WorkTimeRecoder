package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func GetUser(db *xorm.Engine, hash string) (table.Users, error) {
	var user table.Users

	_, err := db.Where(
		"disabled = false",
	).And(
		"hash = ?",
		hash,
	).And(
		"disabled = false",
	).Get(&user)

	return user, err
}

func GetWorkTime(db *xorm.Engine, content string, userId uint64) (table.WorkTimes, error) {
	var workTime table.WorkTimes

	_, err := db.Where(
		"disabled = false",
	).And(
		"content = ?",
		content,
	).And(
		"disabled = false",
	).Get(&workTime)

	return workTime, err
}
