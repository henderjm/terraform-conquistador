package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	AccessKey     string
	SecretKey     string
	CredsFilename string
	Region        string

	//terraformVersion string
}

type AWSClient struct {
	ec2conn   *ec2.EC2
	r53conn   *route53.Route53
	s3conn    *s3.S3
	elbv2conn *elbv2.ELBV2
	region    string
}

func (c *Config) Client() (*AWSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return &AWSClient{}, err
	}

	client := &AWSClient{
		ec2conn:   ec2.New(sess.Copy(&aws.Config{})),
		r53conn:   route53.New(sess.Copy(&aws.Config{})),
		s3conn:    s3.New(sess.Copy(&aws.Config{})),
		elbv2conn: elbv2.New(sess.Copy(&aws.Config{})),
		region:    c.Region,
	}

	return client, nil
}
