package main

import (
	"fmt"
	"time"

	"github.com/byuoitav/mayday/log"

	"github.com/spf13/pflag"
)

func main() {
	alertManager := &AlertManager{
		inDistress: false,
	}
	var port int

	pflag.IntVarP(&port, "port", "p", 9000, "port to host status endpoint")
	pflag.IntVarP(&alertManager.limit, "limit", "l", 100, "arbitrary issue limit")
	pflag.StringVarP(&alertManager.webhook, "webhook", "w", "none", "slack webhook url")
	pflag.Parse()

	if alertManager.webhook == "none" {
		log.P.Panic("Slack webhook needed. Use '-w' at execution.")
	}

	go serveStatus(port)

	for {
		numIssues, err := alertManager.getIssueCount()

		if err != nil {
			log.P.Error(err.Error())
		} else if alertManager.checkLimit(numIssues) {
			err = alertManager.sendAlert(fmt.Sprintf("<!channel> SMEE has reached a total of %d issues", numIssues))
			if err != nil {
				log.P.Error(err.Error())
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
