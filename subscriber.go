package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func subscriber() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	queueUrl := "https://sqs.us-east-2.amazonaws.com/258704584126/testqueue"

	subscribe(queueUrl, sigs)

}

func subscribe(queueUrl string, cancel <-chan os.Signal) {
	awsSession := BuildSession()
	svc := sqs.New(awsSession, nil)

	for {
		messages := receiveMessages(svc, queueUrl)

		for _, msg := range messages {
			if msg == nil {
				continue
			}
			fmt.Println(*msg.Body)
			go DeleteMessage(svc, queueUrl, msg.ReceiptHandle)
		}

		select {
		case <-cancel:
			return
		case <-time.After(1 * time.Millisecond):
		}
	}
}

func receiveMessages(svc *sqs.SQS, queueUrl string) []*sqs.Message {

	receiveMessagesInput := &sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: aws.Int64(10), // max 10
		WaitTimeSeconds:     aws.Int64(1),  // max 20
		VisibilityTimeout:   aws.Int64(1),  // max 20
	}

	receiveMessageOutput, err :=
		svc.ReceiveMessage(receiveMessagesInput)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	if receiveMessageOutput == nil || len(receiveMessageOutput.Messages) == 0 {
		return nil
	}

	return receiveMessageOutput.Messages
}

func DeleteMessage(svc *sqs.SQS, queueUrl string, handle *string) {
	delInput := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: handle,
	}
	_, err := svc.DeleteMessage(delInput)

	if err != nil {
		fmt.Println("Delete Error", err)
		return
	}
}

func SubscribeSNS(session *session.Session, topic string) {
	svc := sns.New(session)

	_, err := svc.Subscribe(&sns.SubscribeInput{
		// Attributes:            nil,
		Endpoint: aws.String("myname@mydomain.com"),
		Protocol: aws.String("email"),
		// ReturnSubscriptionArn: nil,
		TopicArn: aws.String(topic),
	})
	if err != nil {
		fmt.Println(err)
	}
}
