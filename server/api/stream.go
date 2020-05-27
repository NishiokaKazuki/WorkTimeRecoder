package api

import (
	"log"
	"server/config"
	"server/model/db"
	"server/model/enums"
	"server/query"
	"server/utils"
	"time"

	"github.com/slack-go/slack"
)

type param struct {
	date     time.Time
	reportCh string
}

func (p *param) DayReport() (string, error) {
	var (
		ids []uint64
	)
	con := db.GetDBConn()

	workTimes, err := query.FindWorkTimesByDate(con, enums.AllUser, p.date)
	if err != nil {
		return "", err
	}

	for _, w := range workTimes {
		ids = append(ids, w.Id)
	}
	workRests, err := query.FindWorkRestsByDate(con, ids, p.date)
	if err != nil {
		return "", err
	}

	sumTimes, err := utils.CalcWorkTimes(workTimes, workRests, p.date)
	if err != nil {
		return "", err
	}

	return utils.WorkTimeMessage(sumTimes, p.date), nil
}

func StreamServe() {
	var (
		msg string
	)
	log.Println("Starting Stream")

	p := param{
		date: time.Now(),
	}

	conf, err := config.ReadBOTConfig()
	if err != nil {
		return
	}

	rtm := slack.New(conf.TokenID).NewRTM()
	go rtm.ManageConnection()

	tickChan := time.NewTicker(time.Second * 5).C

	for {
		select {
		case <-tickChan:
			if y, m, d := p.date.Date(); utils.DateChange(y, m, d) {
				msg, err = p.DayReport()
				if err != nil {
					msg = err.Error()
				}
				p.date = time.Now()
			}
			rtm.SendMessage(rtm.NewOutgoingMessage(msg, p.reportCh))
		}
	}
}
