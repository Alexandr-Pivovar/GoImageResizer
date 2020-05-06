package infrastrature

import (
	"GoImageZip/internal/app"
	"GoImageZip/internal/domain"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
)

type AWSConnector struct {
	awsSession *s3.S3
	bucket     *string
	url        string
}

// GetImage requests data by url
func (AWSConnector) GetImage(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", app.ErrCould, err)
	}
	defer func() {
		err := r.Body.Close()
		log.Error(err)
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", app.ErrCould, err)
	}

	return b, nil
}

// Save saves input file to aws and return url for to download this file
func (a AWSConnector) Save(id string, image domain.Image) (string, error) {
	_, err := a.awsSession.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image.Data),
		Bucket: a.bucket,
		Key:    aws.String(id),
	})
	if err != nil {
		return "", fmt.Errorf("%s: %s", app.ErrCould, err)
	}

	return a.url + "/" + id + ".png", nil
}

// NewAWSConnector is constructor, receives aws session, bucket and aws url,
//// calls Panic when one of the param not valid
func NewAWSConnector(awsSession *s3.S3, bucket string, url string) *AWSConnector {
	return &AWSConnector{
		awsSession: awsSession,
		bucket:     aws.String(bucket),
		url:        fmt.Sprintf(`https://%s.%s`, bucket, url),
	}
}
