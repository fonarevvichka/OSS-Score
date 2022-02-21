package main

import (
	"api/util"
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

type SQSEvent struct {
	Records []SQSMessage `json:"Records"`
}

type SQSMessage struct {
	MessageId              string            `json:"messageId"` //nolint: stylecheck
	ReceiptHandle          string            `json:"receiptHandle"`
	Body                   string            `json:"body"`
	Md5OfBody              string            `json:"md5OfBody"`
	Md5OfMessageAttributes string            `json:"md5OfMessageAttributes"`
	Attributes             map[string]string `json:"attributes"`
	MessageAttributes      map[string]string `json:"messageAttributes"`
	EventSourceARN         string            `json:"eventSourceARN"`
	EventSource            string            `json:"eventSource"`
	AWSRegion              string            `json:"awsRegion"`
}

type SQSMessageAttribute struct {
	StringValue      *string  `json:"stringValue,omitempty"`
	BinaryValue      []byte   `json:"binaryValue,omitempty"`
	StringListValues []string `json:"stringListValues"`
	BinaryListValues [][]byte `json:"binaryListValues"`
	DataType         string   `json:"dataType"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		if err != nil {
			log.Fatalln("error converting time frame to int")
		}

		util.QueryProject(catalog, owner, name, timeFrame)
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
