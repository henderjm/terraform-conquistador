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
	routeTable routeTable
	elb        elb
}

func NewClient(eN, p, t, oF string) client {
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
	return client{
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
	err = n.Import(c)
	if err != nil {
		return err
	}

	// import route tables
	rt := NewRouteTable()
	err = rt.Import(c)
	if err != nil {
		return err
	}

	// import loadbalancing resources
	lb := NewELB()
	err = lb.Import(c)
	if err != nil {
		return err
	}

	// next resource
	c.awsResources = AwsResources{
		networking: n,
		routeTable: rt,
		elb:        lb,
	}
	return nil
}

func (c *client) UpdateTerraformStateFile() {
	// print networking
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' aws_vpc.base_vpc %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.networking.Vpc.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' aws_internet_gateway.ig %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.networking.Ig.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' aws_route_table.internet_access_through_ig %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.routeTable.InternetGateway.Id))

	// print ALB
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' aws_lb.portal_lb %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.arn.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' 'aws_subnet.lb_subnets[0]' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.subnets[0].Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='%s-vars.tfvars' 'aws_subnet.lb_subnets[1]' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.subnets[1].Id))

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
