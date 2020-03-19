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
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func Working(userName string, content, supplement string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, userName)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	affected, err := query.CreateWorkTime(con, user.Id, content, supplement)
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func FinishWorking(userName string, content string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, userName)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	affected, err := query.UpdateWorkTime(con, content, user.Id)
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success updated, but out range values")
	}

	return nil
}

func Resting(userName, content string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, userName)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	workTime, err := query.GetWorkTime(con, content, user.Id)
	if workTime.Id == 0 {
		return errors.New("Not found worktime. Did you started working?")
	}
	if err != nil {
		return err
	}

	affected, err := query.CreateWorkRest(con, workTime.Id)
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func FinishResting(userName string, content string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, userName)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	workTime, err := query.GetWorkTime(con, content, user.Id)
	if workTime.Id == 0 {
		return errors.New("Not found worktime. Did you started working?")
	}
	if err != nil {
		return err
	}

	affected, err := query.UpdateWorkRest(con, workTime.Id)
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success updated, but out range values")
	}

	return nil
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
