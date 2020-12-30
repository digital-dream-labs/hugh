package s3

import (
	"github.com/aalpern/go-metrics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSSuite struct {
	Suite
	SNS        *sns.SNS
	Config     *aws.Config
	Endpoint   string
	MetricsReg *metrics.Registry
}

func (s *SNSSuite) SetupSuite() {
	s.Start = StartSNSContainer
	s.Suite.SetupSuite()
	s.Endpoint = "http://" + s.Container.Addr()
	s.Config = &aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String(s.Endpoint),
		DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewStaticCredentials("x", "x", ""),
	}
	s.SNS = sns.New(session.New(), s.Config)
}

func (s *SNSSuite) CreateTopic(name string) string {
	topicIn := &sns.CreateTopicInput{
		Name: aws.String(name),
	}
	topicOut, err := s.SNS.CreateTopic(topicIn)
	if err != nil {
		s.T().Fatalf("Failed to create SNS queue: %s", err)
	}
	return *topicOut.TopicArn
}

func (s *SNSSuite) DeleteTopic(arn string) {
	topicIn := &sns.DeleteTopicInput{
		TopicArn: aws.String(arn),
	}
	_, err := s.SNS.DeleteTopic(topicIn)
	if err != nil {
		s.T().Fatalf("Failed to delete SNS topic:: %s", err)
	}
	return
}
