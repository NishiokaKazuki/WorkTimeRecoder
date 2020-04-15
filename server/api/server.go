package api

import (
	"errors"
	"log"
	"server/config"
	"server/model/db"
	"server/model/table"
	"server/query"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type Slackparams struct {
	tokenID   string
	botID     string
	workingCh string
	reportCh  string
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
		affected, err := query.InsertUser(con, appUser.RealName, appUser.ID)
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

	affected, err := query.InsertWorkTime(
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

	affected, err := query.UpdateWorkTime(con,
		table.WorkTimes{
			UserId:     user.Id,
			Content:    message[0],
			FinishedAt: date,
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

	affected, err := query.InsertWorkRest(con, table.WorkRests{
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

func FinishResting(hash string, message []string) error {
	date := time.Now()
	con := db.GetDBConn()

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

	affected, err := query.UpdateWorkRest(con, table.WorkRests{
		WorkTimeId: workTime.Id,
		FinishedAt: date,
	})
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success updated, but out range values")
	}

	return nil
}

func SuspensionWorking(hash string, message []string) (string, error) {
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return "", errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return "", err
	}

	workTimeId, err := strconv.ParseUint(message[0], 10, 64)
	if err != nil {
		return "", err
	}

	affected, err := query.DeleteWorkTimes(con, workTimeId, user.Id)
	if err != nil {
		return "", err
	}
	if affected != true {
		return "", errors.New("Success deleted, but out range values")
	}

	return "", nil
}

func WorkLog(hash string, message []string) (string, error) {
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return "", errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return "", err
	}

	cnt, err := strconv.Atoi(message[0])
	if err != nil {
		return "", err
	}

	workTimes, err := query.FindWorkTimeLatest(con, cnt, user.Id)
	if err != nil {
		return "", err
	}

	return utils.WorkLogMessage(workTimes)
}

func WorkTime(hash string, message []string) (string, error) {
	date := time.Now()
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return "", errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return "", err
	}

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	workTimes, err := query.FindWorkTimesByDate(con, user.Id, date)
	if err != nil {
		return "", err
	}

	workRests, err := query.FindWorkRestsByDate(con, user.Id, date)
	if err != nil {
		return "", err
	}

	sumTimes, err := utils.CalcWorkTimes(workTimes, workRests, date)
	if err != nil {
		return "", err
	}

	return utils.WorkTimeMessage(sumTimes, date), nil
}

func WorkInfo(hash string, message []string) (string, error) {
	date := time.Now()
	con := db.GetDBConn()

	user, err := query.GetUser(con, hash)
	if user.Id == 0 {
		return "", errors.New("Not found user. Did you completed SignUp?")
	}
	if err != nil {
		return "", err
	}

	workInfo, err := query.FindWorkInfos(con, date, user.Id)
	if err != nil {
		return "", err
	}

	return utils.SplitWorkInfo(workInfo, user)
}

func (s *Slackparams) ValidateMessageEvent(ev *slack.MessageEvent) error {

	switch ev.Channel {
	case s.signupCh:
		user, err := s.rtm.GetUserInfo(ev.Msg.User)
		if err != nil {
			return err
		}
		if err := SignUp(user); err != nil {
			return err
		}
	case s.workingCh:
		if strings.HasPrefix(ev.Msg.Text, s.botID) {
			res, err := PrefixMessage(ev.Msg.Text)
			if err != nil {
				return err
			}
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		} else {
			res, err := WorkingMessage(ev.Msg.User, ev.Msg.Text)
			if err != nil {
				s.rtm.SendMessage(s.rtm.NewOutgoingMessage(err.Error(), ev.Channel))
				return err
			}
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		}
	case s.reportCh:
		if strings.HasPrefix(ev.Msg.Text, s.botID) {
			res, err := PrefixMessage(ev.Msg.Text)
			if err != nil {
				return err
			}
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		} else {
			res, err := ReportMessage(ev.Msg.User, ev.Msg.Text)
			if err != nil {
				s.rtm.SendMessage(s.rtm.NewOutgoingMessage(err.Error(), ev.Channel))
				return err
			}
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
		}
	default:
		log.Println("%s %s", ev.Channel, ev.Msg.Text)
	}

	return nil
}

func PrefixMessage(message string) (string, error) {
	var res string

	log.Println(message)
	switch message {
	default:
		res = "やあ僕きもかわいいgopher君!僕は以下の機能を持ってるよ活用してね!\n\n" +
			"作業記録機能(例)\n" +
			"1.開始 `hogehoge` -tm\n" +
			"2.中断 `hogehoge` -t\n" +
			"3.再開 `hogehoge` -t\n" +
			"4.終了 `hogehoge` -t\n\n" +
			"option list\n" +
			"-t  : 時間を指定可能 `yyyy-mm-dd hh:mm`, `hh:mm` のフォーマット\n" +
			"-m  : コメントを追記可能 ただし半角空白で区切られる\n"
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
		if err := FinishResting(hash, m[1:]); err != nil {
			return "", err
		}
		res = "End Resting"
	default:
		res = "no response"
	}
	return res, nil
}

func ReportMessage(hash, message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")

	switch m[0] {
	case "作業記録":
		r, err := WorkInfo(hash, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "作業時間":
		r, err := WorkTime(hash, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "log":
		r, err := WorkLog(hash, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "rm":
		r, err := SuspensionWorking(hash, m[1:])
		if err != nil {
			return "", err
		}
		res = r
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
		workingCh: conf.WorkingCh,
		reportCh:  conf.ReportCh,
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
