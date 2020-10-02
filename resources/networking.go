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
	Vpc     AWSResourceId
	Ig      AWSResourceId
	Subnets []AWSResourceId
}

func NewNetworking() networking { return networking{} }

func (n *networking) Import(c *client) (networking, error) {
	vpc, err := importVPC(c)
	if err != nil {
		return networking{}, err
	}
	ig, err := importIg(c)
	if err != nil {
		return networking{}, err
	}
	subnets, err := importSubnets(c)
	if err != nil {
		return networking{}, err
	}
	return networking{
		Vpc:     vpc,
		Ig:      ig,
		Subnets: subnets,
	}, nil
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

func importSubnets(c *client) ([]AWSResourceId, error) {
	fmt.Println("searching-for-subnets")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
		},
	}

	result, err := conn.DescribeSubnets(input)
	if err != nil {
		handleAWSError(err)
		return []AWSResourceId{}, err
	}

	var subnets []AWSResourceId

	for _, s := range result.Subnets {
		r := AWSResourceId{
			Id: s.SubnetId,
		}
		subnets = append(subnets, r)

	}

	return subnets, nil
}
