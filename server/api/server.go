package api

import (
	"errors"
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

	log.Println("check")
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	if strings.HasPrefix(ev.Msg.Text, s.botID) {
		res, err := PrefixMessage(ev.Msg.Text)
		if err != nil {
			return err
		}
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		return nil
	} else {
		res, err := WorkingMessage(ev.Msg.Text)
		if err != nil {
			return err
		}
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		return nil
	}

	return nil
}

func PrefixMessage(message string) (string, error) {
	var res string

	log.Println(message)
	switch message {
	default:
		res = "Command List\n"
	}

	return res, nil
}

func WorkingMessage(message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")

	switch m[0] {
	case "開始":
		res = "Start Working"
	case "終了":
		res = "End Working"
	case "中断":
		res = "Start Resting"
	case "再開":
		res = "End Resting"
	default:
		res = "no response"
	}
	return res, nil
}

func ListenAndServe(token string) {
	log.Println("Starting Server")

	params := Slackparams{
		tokenID:   "xoxb-1015900425939-1027655437248-OXFLwNFfN7UwMl9Q7ItalPrK",
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
				log.Println("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}
