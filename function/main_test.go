package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

var (
	exampleSNS = `
  {
    "Records": [
      {
        "EventSource": "aws:sns",
        "EventVersion": "1.0",
        "EventSubscriptionArn": "arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "Sns": {
          "Type": "Notification",
          "MessageId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
          "TopicArn": "arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms",
          "Subject": "ALARM: \"Example alarm name\" in EU - Ireland",
          "Message": "{\"AlarmName\":\"Example alarm name\",\"AlarmDescription\":\"Example alarm description.\",\"AWSAccountId\":\"000000000000\",\"NewStateValue\":\"ALARM\",\"NewStateReason\":\"Threshold Crossed: 1 datapoint (10.0) was greater than or equal to the threshold (1.0).\",\"StateChangeTime\":\"2019-01-12T15:30:42.136+0000\",\"Region\":\"EU - Ireland\",\"OldStateValue\":\"OK\",\"Trigger\":{\"MetricName\":\"DeliveryErrors\",\"Namespace\":\"ExampleNamespace\",\"Statistic\":\"SUM\",\"Unit\":null,\"Dimensions\":[],\"Period\":300,\"EvaluationPeriods\":1,\"ComparisonOperator\":\"GreaterThanOrEqualToThreshold\",\"Threshold\":1.0}}",
          "Timestamp": "2019-01-12T15:30:42.119Z",
          "SignatureVersion": "1",
          "Signature": "AA==",
          "SigningCertUrl": "https://sns.eu-west-1.amazonaws.com/SimpleNotificationService-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.pem",
          "UnsubscribeUrl": "https://sns.eu-west-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
          "MessageAttributes": {}
        }
      }
    ]
  }
`
)

func TestHandler(t *testing.T) {
	var sns events.SNSEvent
	err := json.Unmarshal([]byte(exampleSNS), &sns)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	handleRequest(sns)
}
