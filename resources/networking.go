package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/*
VPC
Internet Gateway

*/
type networking struct {
	Vpc AWSResourceId
	Ig  AWSResourceId
}

func NewNetworking() networking { return networking{} }

func (net *networking) Import(c *client) error {
	var err error
	net.Vpc, err = importVPC(c)
	if err != nil {
		return err
	}
	net.Ig, err = importIg(c)
	if err != nil {
		return err
	}
	return nil
}

func importVPC(c *client) (AWSResourceId, error) {
	fmt.Println("searching-for-vpc")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
		},
	}

	result, err := conn.DescribeVpcs(input)
	if err != nil {
		handleAWSError(err)
		return AWSResourceId{}, err
	}

	if len(result.Vpcs) != 1 {
		return AWSResourceId{}, errors.New(fmt.Sprintf("found: %d vpc(s), should only find 1", len(result.Vpcs)))
	}

	vpc := AWSResourceId{
		Id: result.Vpcs[0].VpcId,
	}

	return vpc, nil
}

func importIg(c *client) (AWSResourceId, error) {
	fmt.Println("searching-for-ig")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
		},
	}

	result, err := conn.DescribeInternetGateways(input)
	if err != nil {
		handleAWSError(err)
		return AWSResourceId{}, err
	}

	if len(result.InternetGateways) != 1 {
		return AWSResourceId{}, errors.New(fmt.Sprintf("found: %d ig(s), should only find 1", len(result.InternetGateways)))
	}

	ig := AWSResourceId{
		Id: result.InternetGateways[0].InternetGatewayId,
	}

	return ig, nil
}
