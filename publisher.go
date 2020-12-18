package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func publisher() {

	sigs := make(chan os.Signal, 1) // make channel
	
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	destination := "https://sqs.us-east-2.amazonaws.com/258704584126/testqueue" // SQS url 

	go publishMessagesFromStdIn(SendSQS, destination)

	<-sigs

}

func publishMessagesFromStdIn(sender func(session *session.Session, destination string, message string), destination string) {
	awsSession := BuildSession()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}

		sender(awsSession, destination, text[:len(text)-1])
	}
}

func SendSNS(session *session.Session, destination string, message string) {
	svc := sns.New(session)

	pubInput := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(destination),
	}

	_, err := svc.Publish(pubInput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fmt.Println(output.MessageId)
}

func SendSQS(session *session.Session, destination string, message string) {
	svc := sqs.New(session, nil)

	sendInput := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(destination),
	}

	_, err := svc.SendMessage(sendInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(output.MessageId)
}
