package api

import (
	"errors"
	"fmt"
	"log"
	"server/model/db"
	"server/query"
	"strings"

	"github.com/slack-go/slack"
)

type Slackparams struct {
	tokenID   string
	botID     string
	channelID string
	rtm       *slack.RTM
}

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

func (s *Slackparams) ValidateMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, s.botID) {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Parse message start
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 {
		return fmt.Errorf("invalid message")
	}

	if m[0] == "ホリネズミ?" {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("そうだ！", ev.Channel))
		return nil
	}

	if m[0] == "ネズミ?" {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("ちがう！", ev.Channel))
		return nil
	}

	return nil
}

func ListenAndServe(token string) {
	log.Println("Starting Server")

	params := Slackparams{
		tokenID:   "xoxb-1015900425939-1027655437248-ZIst9HIL1KF8z89ThuUL3Hce",
		botID:     "<@U010TK9CV7A>",
		channelID: "C010HHPLTFB",
	}

	api := slack.New(params.tokenID)

	params.rtm = api.NewRTM()
	go params.rtm.ManageConnection()

	for msg := range params.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := params.ValidateMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}
