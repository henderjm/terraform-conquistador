package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"

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
	ec2conn *ec2.EC2
	r53conn *route53.Route53
	region  string
	s3conn  *s3.S3
}

func (c *Config) Client() (*AWSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return &AWSClient{}, err
	}

	//awsbaseConfig := &awsbase.Config{
	//	AccessKey: c.AccessKey,
	//	Region:    c.Region,
	//	SecretKey: c.SecretKey,
	//}

	//sess, accountID, partition, err := awsbase.GetSessionWithAccountIDAndPartition(awsbaseConfig)
	//if err != nil {
	//	return nil, fmt.Errorf("error configuring Terraform AWS Provider: %w", err)
	//}

	client := &AWSClient{
		ec2conn: ec2.New(sess.Copy(&aws.Config{})),
		region:  c.Region,
	}

	route53Config := &aws.Config{}

	client.s3conn = s3.New(sess.Copy(&aws.Config{}))

	client.r53conn = route53.New(sess.Copy(route53Config))

	return client, nil
}
