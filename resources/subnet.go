package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/*
aws routeTable
*/

type subnet struct {
	Subnets []AWSResourceId
}

func NewSubnet() subnet { return subnet{} }

func (s *subnet) Import(c *client) (subnet, error) {
	subnets, err := importSubnets(c)
	if err != nil {
		return subnet{}, err
	}
	return subnet{
		Subnets: subnets,
	}, nil
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
