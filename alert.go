package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mayday/log"
	"net/http"

	"github.com/slack-go/slack"
)

type AlertManager struct {
	inDistress bool
	limit      int
	webhook    string
}

func (am *AlertManager) getIssueCount() (int, error) {
	url := "https://smee.av.byu.edu/issues"
	log.P.Info(fmt.Sprintf("sending request to: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var issues []interface{}
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return -1, err
	}

	log.P.Info(fmt.Sprintf("received %d issues", len(issues)))
	return len(issues), nil
}

func (am *AlertManager) checkLimit(issues int) bool {
	if issues >= am.limit && !am.inDistress {
		am.inDistress = true
		return true
	} else if issues < am.limit {
		am.inDistress = false
	}
	return false
}

func (am *AlertManager) sendAlert(msg string) error {
	log.P.Info(fmt.Sprintf("sending alert to slack: %s", msg))
	alertAttachment := slack.Attachment{
		Text: msg,
	}

	slackMessage := slack.WebhookMessage{
		Attachments: []slack.Attachment{alertAttachment},
	}

	return slack.PostWebhook(am.webhook, &slackMessage)
}
