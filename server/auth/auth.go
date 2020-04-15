package auth

import (
	"server/model/db"
	"server/model/table"
	"server/query"
)

func Auth(token string) (bool, table.Users) {
	var user table.Users
	con := db.GetDBConn()

	user, err := query.GetUser(con, token)
	if user.Id == 0 {
		return false, user
	}
	if err != nil {
		return false, user
	}

	return true, user
}
