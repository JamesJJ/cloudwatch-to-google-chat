package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	SuppressableAlertsStateFromTo = map[string]string{
		"INSUFFICIENT_DATA": "OK",
	}
)

func handleRequest(request events.SNSEvent) error {

	var cloudwatchAlarm events.CloudWatchAlarmSNSPayload
	err := json.Unmarshal([]byte(request.Records[0].SNS.Message), &cloudwatchAlarm)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return nil // Non-retryable error
	}

	if suppress, err := strconv.ParseBool(os.Getenv("SUPPRESS_UNKNOWN_TO_OK")); err == nil && suppress {
		if SuppressableAlertsStateFromTo[cloudwatchAlarm.OldStateValue] == cloudwatchAlarm.NewStateValue {
			log.Printf(
				"Alarm suppressed for %s --> %s \"%s\"",
				cloudwatchAlarm.OldStateValue,
				cloudwatchAlarm.NewStateValue,
				cloudwatchAlarm.AlarmName,
			)
			return nil
		}
	}

	log.Printf(
		"Alarm: %s --> %s \"%s\"",
		cloudwatchAlarm.OldStateValue,
		cloudwatchAlarm.NewStateValue,
		cloudwatchAlarm.AlarmName,
	)

	postUrl, ok := os.LookupEnv("CHAT_WEBHOOK")
	if !ok {
		log.Printf("Must set CHAT_WEBHOOK environment variable to the web hook URL")
		return nil // Non-retryable error
	}

	chatMessage := messageFromAlarm(cloudwatchAlarm)
	if err := postMessageToChat(postUrl, chatMessage); err != nil {
		return err // Retryable
	}
	log.Println("Message sent")
	return nil
}

func messageFromAlarm(alarm events.CloudWatchAlarmSNSPayload) Message {

	alertTimeString := html.EscapeString(alarm.StateChangeTime)
	if alertTime, err := time.Parse("2006-01-02T15:04:05.000-0700", alarm.StateChangeTime); err == nil {
		alertTimeString = alertTime.Format("Monday, 02 Jan 2006 15:04:05 -0700")
	}

	msg := Message{
		Cards: []MessageCard{
			{
				CardHeader: CardHeader{
					Title: alarm.AlarmName,
					Subtitle: fmt.Sprintf(
						"Cloudwatch: %s",
						alarm.NewStateValue,
					),
				},
				Sections: []CardSection{
					{
						Widgets: []Widget{
							{
								TextParagraph: &TextParagraph{
									Text: html.EscapeString(alarm.AlarmDescription),
								},
							},
							{
								KeyValue: &WidgetKeyValue{
									TopLabel: "Transition",
									Content: fmt.Sprintf(
										"%s &mdash;&mdash;&gt; %s",
										alarm.OldStateValue,
										alarm.NewStateValue,
									),
								},
							},
							{
								TextParagraph: &TextParagraph{
									Text: fmt.Sprintf(
										"<i>%s</i>",
										html.EscapeString(alarm.NewStateReason),
									),
								},
							},
							{
								KeyValue: &WidgetKeyValue{
									TopLabel: "Source",
									Content: fmt.Sprintf(
										"%s / %s",
										html.EscapeString(alarm.AWSAccountID),
										html.EscapeString(alarm.Region),
									),
								},
							},
							{
								KeyValue: &WidgetKeyValue{
									TopLabel: "Time",
									Content:  alertTimeString,
								},
							},
						},
					},
				},
			},
		},
	}

	switch alarm.NewStateValue {
	case "OK":
		if imageUrl := os.Getenv("OK_IMAGE_URL"); imageUrl != "" {
			for i := range msg.Cards {
				msg.Cards[i].CardHeader.ImageUrl = imageUrl
			}
		}
	default:
		if imageUrl := os.Getenv("ALERT_IMAGE_URL"); imageUrl != "" {
			for i := range msg.Cards {
				msg.Cards[i].CardHeader.ImageUrl = imageUrl
			}
		}
	}
	return msg
}

func postMessageToChat(postUrl string, message Message) error {
	client := &http.Client{}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// log.Println(string(data))
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("POST response status: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	runtime.Start(handleRequest)
}
