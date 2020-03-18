package api

import (
	"errors"
	"server/model/db"
	"server/query"

	"github.com/nlopes/slack"
)

func SignUp(userName string) error {
	con := db.GetDBConn()

	affected, err := query.CreateUser(con, userName)
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return err
}

func Working(userName string, content, supplement string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, userName)
	if err != nil {
		return err
	}

	affected, err := query.CreateWorkTime(con, user.Id, content, supplement)
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return err
}

func ListenAndServe(token string) {
	api := slack.New(
		token,
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			rtm.SendMessage(rtm.NewOutgoingMessage("ホリネズミです。塊茎(球根)食べたい", ev.Channel))
		}
	}
}
