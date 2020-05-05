package main

import (
	"GoImageZip/internal/app"
	"GoImageZip/internal/app/mocks"
	"GoImageZip/internal/infrastrature"
	"GoImageZip/internal/interfaces"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
	"log"
)

type AwsConf struct {
	Url    string
	Bucket string
	Region string
	Id     string
	Secret string
	Token  string
}

func main() {
	addr := flag.String("a", "localhost:8080", "-a <addr:port/>")
	redisAddr := flag.String("r", "localhost:6379", "-r <addr:port/>")
	redisPass := flag.String("p", "", "-p <redis password>")
	redisDB := flag.Int("d", 0, "-d <redis DB number>")
	awsID := flag.String("w", "", "-w <aws id number>")
	awsSecret := flag.String("s", "", "-s <aws secret>")
	awsToken := flag.String("t", "", "-t <aws token>")
	awsEndPoint := flag.String("e", "", "-e <aws endpoint>")
	awsBucket := flag.String("b", "", "-b <aws bucket>")
	awsRegeon := flag.String("g", "", "-g <aws region>")

	flag.Parse()

	newAWSSession := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*awsID, *awsSecret, *awsToken),
		Endpoint:    aws.String("https://" + *awsEndPoint),
		Region:      aws.String(*awsRegeon),
	})

	infrastrature.NewAWSConnector(s3.New(newAWSSession), *awsBucket, *awsEndPoint)

	redisConn, err := infrastrature.NewRedisConnector(*redisAddr, *redisPass, *redisDB)
	if err != nil {
		log.Fatalln(err)
	}
	repo := interfaces.NewRedisRepo(redisConn)

	service := app.NewImageService(repo, &interfaces.ImageResize{}, func() *mocks.Clouder {
		m := &mocks.Clouder{}
		m.On("Save", mock.Anything, mock.Anything).Return("http://domain", nil)
		return m
	}())

	interfaces.NewController(service).Run(*addr)
}
