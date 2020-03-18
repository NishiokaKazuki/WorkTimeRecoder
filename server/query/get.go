package query

import (
	"server/model/table"

	"github.com/go-xorm/xorm"
)

func GetUser(db *xorm.Engine, name string) (table.Users, error) {
	var user table.Users

	_, err := db.Where(
		"disabled = false",
	).And(
		"name = ?",
		name,
	).Get(&user)

	return user, err
}
