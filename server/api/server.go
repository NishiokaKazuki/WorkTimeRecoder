package api

import (
	"errors"
	"log"
	"server/auth"
	"server/config"
	"server/model/db"
	"server/model/sentence"
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

func Working(user table.Users, message []string) error {
	var (
		affected bool
		err      error
	)
	supplement := ""
	content := message[0]
	date := time.Now()
	con := db.GetDBConn()

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	if supple, has := utils.SplitSuppleOption(message[1:]); has == true {
		supplement = supple
	}

	workTime, _ := query.GetWorkTime(con, content, user.Id)
	if workTime.Id != 0 {
		affected, err = query.UpdateStartOnWorkTime(
			con,
			table.WorkTimes{
				UserId:     user.Id,
				Content:    content,
				Supplement: supplement,
				StartedAt:  date,
			},
		)
	} else {
		affected, err = query.InsertWorkTime(
			con,
			table.WorkTimes{
				UserId:     user.Id,
				Content:    content,
				Supplement: supplement,
				StartedAt:  date,
			},
		)
	}
	if err != nil {
		return err
	}
	if affected != true {
		return errors.New("Success created, but out range values")
	}

	return nil
}

func FinishWorking(user table.Users, message []string) error {
	date := time.Now()
	con := db.GetDBConn()

	if d, has := utils.SplitTimeOption(message[1:]); has == true {
		date = d
	}

	affected, err := query.UpdateWorkTime(con,
		table.WorkTimes{
			UserId:     user.Id,
			Content:    message[0],
			IsFinished: true,
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

func Resting(user table.Users, message []string) error {
	date := time.Now()
	con := db.GetDBConn()

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

func FinishResting(user table.Users, message []string) error {
	date := time.Now()
	con := db.GetDBConn()

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
		IsFinished: true,
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

func SuspensionWorking(user table.Users, message []string) (string, error) {
	con := db.GetDBConn()

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

func WorkLog(user table.Users, message []string) (string, error) {
	con := db.GetDBConn()

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

func WorkTime(user table.Users, message []string) (string, error) {
	var ids []uint64
	date := time.Now()
	con := db.GetDBConn()

	if d, has := utils.SplitTimeOption(message[0:]); has == true {
		date = d
	}

	workTimes, err := query.FindWorkTimesByDate(con, user.Id, date)
	if err != nil {
		return "", err
	}

	for _, w := range workTimes {
		ids = append(ids, w.Id)
	}
	workRests, err := query.FindWorkRestsByDate(con, ids, date)
	if err != nil {
		return "", err
	}

	sumTimes, err := utils.CalcWorkTimes(workTimes, workRests, date)
	if err != nil {
		return "", err
	}

	return utils.WorkTimeMessage(sumTimes, date), nil
}

func WorkInfo(user table.Users, message []string) (string, error) {
	date := time.Now()
	con := db.GetDBConn()

	workInfo, err := query.FindWorkInfos(con, date, user.Id)
	if err != nil {
		return "", err
	}

	return utils.SplitWorkInfo(workInfo, user)
}

func (s *Slackparams) ValidateMessageEvent(ev *slack.MessageEvent) error {
	var err error
	res := ""

	u, _ := s.rtm.GetUserInfo(ev.Msg.User)
	has, user := auth.Auth(u.ID)

	switch ev.Channel {
	case s.signupCh:
		err = SignUp(u)

	case s.workingCh:
		if strings.HasPrefix(ev.Msg.Text, s.botID) {
			res, err = PrefixMessage(ev.Msg.Text)
		} else {
			if has != true {
				return errors.New("SignUp still hasn't been completed.")
			}
			res, err = WorkingMessage(user, ev.Msg.Text)
		}

	case s.reportCh:
		if strings.HasPrefix(ev.Msg.Text, s.botID) {
			res, err = PrefixMessage(ev.Msg.Text)
		} else {
			if has != true {
				return errors.New("SignUp still hasn't been completed.")
			}
			res, err = ReportMessage(user, ev.Msg.Text)
		}

	default:
		log.Println("%s %s", ev.Channel, ev.Msg.Text)
	}

	if err != nil {
		res = err.Error()
	}
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(res, ev.Channel))
	return err
}

func PrefixMessage(message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")
	if len(m) <= 1 {
		return sentence.Greeting, nil
	}

	switch m[1] {
	case "fuck":
		res = sentence.Fuck
	case "ブディの真似して":
		res = sentence.Budi
	case "やあ":
		fallthrough
	default:
		res = sentence.Greeting
	}

	return res, nil
}

func WorkingMessage(user table.Users, message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")

	switch m[0] {
	case "開始":
		if err := Working(user, m[1:]); err != nil {
			return "", err
		}
		res = "Start Working"
	case "終了":
		if err := FinishWorking(user, m[1:]); err != nil {
			return "", err
		}
		res = "End Working"
	case "中断":
		if err := Resting(user, m[1:]); err != nil {
			return "", err
		}
		res = "Start Resting"
	case "再開":
		if err := FinishResting(user, m[1:]); err != nil {
			return "", err
		}
		res = "End Resting"
	case "help":
		res = sentence.Help
	default:
		res = "no response"
	}
	return res, nil
}

func ReportMessage(user table.Users, message string) (string, error) {
	var res string

	m := strings.Split(strings.TrimSpace(message), " ")

	switch m[0] {
	case "作業記録":
		r, err := WorkInfo(user, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "作業時間":
		r, err := WorkTime(user, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "log":
		r, err := WorkLog(user, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "rm":
		r, err := SuspensionWorking(user, m[1:])
		if err != nil {
			return "", err
		}
		res = r
	case "help":
		res = sentence.HelpReport
	default:
		res = "no response"
	}

	return res, nil
}

func ListenAndServe(token string) {
	log.Println("Starting Server")
	log.Println(db.GetDBConn().GetTZDatabase())
	log.Println(db.GetDBConn().GetTZLocation())

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
