package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
)

type client struct {
	envName      string
	phase        string
	tags         []string
	outputFile   string
	awsResources AwsResources
	awsClient    *AWSClient
}

type AwsResources struct {
	networking networking
	elb        elb
}

func NewClient(eN, p, t, oF string) *client {
	fmt.Println("creating client")
	tags := strings.Split(t, ",")
	config := Config{
		AccessKey:     "access_key",
		SecretKey:     "secret_key",
		CredsFilename: "~/.aws/credentials",
		Region:        "eu-west-1",
	}
	aC, err := config.Client()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return &client{
		envName:    eN,
		phase:      p,
		tags:       tags,
		outputFile: oF,
		awsClient:  aC,
	}
}

func (c *client) ImportTerraformResources() error {
	var err error
	// import networking
	n := NewNetworking()
	c.awsResources.networking, err = n.Import(c)
	if err != nil {
		return err
	}

	lb := NewELB()
	c.awsResources.elb, err = lb.Import(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) UpdateTerraformStateFile() {
	// print networking
	fmt.Printf("terraform import -var-file=vars.tfvars aws_vpc.base_vpc %s\n", aws.StringValue(c.awsResources.networking.Vpc.Id))
	fmt.Printf("terraform import -var-file=vars.tfvars aws_internet_gateway.ig %s\n", aws.StringValue(c.awsResources.networking.Ig.Id))

	// print ALB
	fmt.Printf("terraform import -var-file=vars.tfvars aws_lb.portal_lb %s\n", aws.StringValue(c.awsResources.elb.alb.Id))

	// print ALB security groups
}

func handleAWSError(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return err
}
