package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	pb "github.com/goniapm/goniplus-worker/metric"
	"log"
	"net/http"
)

const (
	slackQuery   = "select distinct url from notification_slack join project on project.id = notification_slack.project_id where project.apikey = ?"
	titlePanic   = "Panic Alert"
	titleSlow    = "Slow Transaction Alert"
	colorAlert   = "#ff7595"
	colorWarning = "#ffda00"
)

type slackMsgField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type slackMsg struct {
	Fallback  string          `json:"fallback"`
	Color     string          `json:"color"`
	Path      string          `json:"author_name"`
	Fields    []slackMsgField `json:"fields"`
	Instance  string          `json:"footer"`
	Timestamp string          `json:"ts"`
}

type slackBody struct {
	Attachments []slackMsg `json:"attachments"`
}

func createMessage(title, color, path, instance, timestamp string, count int) slackMsg {
	return slackMsg{
		Fallback: fmt.Sprintf("[%s] %s", title, path),
		Color:    color,
		Path:     path,
		Fields: []slackMsgField{
			{
				Title: title,
				Value: fmt.Sprintf("%d transaction(s)", count),
				Short: false,
			},
		},
		Instance:  instance,
		Timestamp: timestamp,
	}
}

func createMessageBody(instance, timestamp string, panicMap, slowMap map[string]int) ([]byte, error) {
	var msg []slackMsg
	for path, count := range panicMap {
		msg = append(msg,
			createMessage(titlePanic, colorAlert, path, instance, timestamp, count))
	}
	for path, count := range slowMap {
		msg = append(msg,
			createMessage(titleSlow, colorWarning, path, instance, timestamp, count))
	}
	body := slackBody{Attachments: msg}
	msgByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return msgByte, nil
}

func sendSlackNotification(metric *pb.Metric) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered at sendSlackNotification: ", r)
		}
	}()
	transactionMetric := metric.Application.Transaction.Detail
	panicMap := make(map[string]int)
	slowMap := make(map[string]int)
	for _, detail := range transactionMetric {
		for _, transaction := range detail.Data {
			if transaction.Status.Panic {
				panicMap[detail.Path]++
			}
			if transaction.Status.Duration > 3000 {
				slowMap[detail.Path]++
			}
		}
	}
	if len(panicMap) == 0 && len(slowMap) == 0 {
		return
	}
	body, err := createMessageBody(metric.Instance, metric.Timestamp, panicMap, slowMap)
	if err != nil {
		log.Println(err)
		return
	}
	if err := mySQL.Ping(); err != nil {
		log.Println(err)
		return
	}
	var url string
	err = mySQL.QueryRow(slackQuery, metric.Apikey).Scan(&url)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	req.Close = true
	res, err := httpClient.Do(req)
	res.Body.Close()
	if err != nil {
		log.Println(err)
	}
}
