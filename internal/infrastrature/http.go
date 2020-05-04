package infrastrature

import (
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

func (AWSConnector) GetImage(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := r.Body.Close()
		log.Error(err)
	}()

	return ioutil.ReadAll(r.Body)
}

// Save saves input file to aws and return url for to download this file
// returns error if PutObject returns error
// Input example: `name:{1.txt},dataUrl:{data:text/plain;base64,MQ==}`
// Where is name - filename with extension, dataUrl - file body in dataURL format
func (a AWSConnector) Save(id string,image domain.Image) (string, error) {
	_, err := a.awsSession.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image.Data),
		Bucket: a.bucket,
		Key:    aws.String(id),
	})

	return a.url + "/" + id + "." + image.Format, err
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

