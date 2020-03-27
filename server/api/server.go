package api

import (
	"errors"
	"log"
	"server/config"
	"server/model/db"
	"server/model/table"
	"server/query"
	"server/utils"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type Slackparams struct {
	tokenID   string
	botID     string
	channelID string
	signupCh  string
	rtm       *slack.RTM
}

func SignUp(appUser *slack.User) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, appUser.ID)
	if err != nil {
		return err
	}
	if user.Id != 0 {
		affected, err := query.UpdateUser(con, appUser.RealName, appUser.ID)
		if err != nil {
			return err
		}
		if affected != true {
			return errors.New("Success updated, but out range values")
		}
	} else {
		affected, err := query.CreateUser(con, appUser.RealName, appUser.ID)
		if err != nil {
			return err
		}
		if affected != true {
			return errors.New("Success created, but out range values")
		}
	}

	return nil
}

func Working(hash string, message []string) error {
	supplement := ""
	date := time.Now()
	con := db.GetDBConn()

	if len(message) == 0 {
		return errors.New("Too few arguments for working.")
	}

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	if supple, has := utils.SplitSuppleOption(message[1:]); has == true {
		supplement = supple
	}

	affected, err := query.CreateWorkTime(
		con,
		table.WorkTimes{
			UserId:     user.Id,
			Content:    message[0],
			Supplement: supplement,
			StartedAt:  date,
		},
	)
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func FinishWorking(hash string, message []string) error {
	supplement := ""
	date := time.Now()
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	if supple, has := utils.SplitSuppleOption(message[1:]); has == true {
		supplement = supple
	}

	affected, err := query.UpdateWorkTime(con,
		table.WorkTimes{
			UserId:     user.Id,
			Content:    message[0],
			Supplement: supplement,
			StartedAt:  date,
		})
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success updated, but out range values")
	}

	return nil
}

func Resting(hash string, message []string) error {
	date := time.Now()
	con := db.GetDBConn()

	if len(message) == 0 {
		return errors.New("Too few arguments for working.")
	}

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return err
	}

	workTime, err := query.GetWorkTime(con, message[0], user.Id)
	if workTime.Id == 0 {
		return errors.New("Not found worktime. Did you started working?")
	}
	if err != nil {
		return err
	}

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	affected, err := query.CreateWorkRest(con, table.WorkRests{
		WorkTimeId: workTime.Id,
		StartedAt:  date,
	})
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func FinishResting(hash string, content string) error {
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
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

	if ev.Channel == s.signupCh {
		user, err := s.rtm.GetUserInfo(ev.Msg.User)
		if err != nil {
			return err
		}
		if err := SignUp(user); err != nil {
			return err
		}
		return nil
	}

	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Println("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	if strings.HasPrefix(ev.Msg.Text, s.botID) {
		res, err := PrefixMessage(ev.Msg.Text)
		if err != nil {
			return err
		}
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
	} else {
		res, err := WorkingMessage(ev.Msg.User, ev.Msg.Text)
		if err != nil {
			return err
		}
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
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

func WorkingMessage(hash, message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")

	switch m[0] {
	case "開始":
		if err := Working(hash, m[1:]); err != nil {
			return "", err
		}
		res = "Start Working"
	case "終了":
		if err := FinishWorking(hash, m[1:]); err != nil {
			return "", err
		}
		res = "End Working"
	case "中断":
		if err := Resting(hash, m[1:]); err != nil {
			return "", err
		}
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

	conf, err := config.ReadBOTConfig()
	if err != nil {
		return
	}
	params := Slackparams{
		tokenID:   conf.TokenID,
		botID:     conf.BotID,
		channelID: conf.ChannelID,
		signupCh:  conf.SignUpCh,
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
