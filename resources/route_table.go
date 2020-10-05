package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/*
aws routeTable
*/

type routeTable struct {
	InternetGateway AWSResourceId
}

func NewRouteTable() routeTable { return routeTable{} }

func (r *routeTable) Import(c *client) error {
	err := r.importRouteTables(c)
	if err != nil {
		return err
	}

	return nil
}

func (r *routeTable) importRouteTables(c *client) error {
	fmt.Println("searching-for-found_rts")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("%s-Pub-RT", c.envName))},
			},
		},
	}

	result, err := conn.DescribeRouteTables(input)
	if err != nil {
		handleAWSError(err)
		return err
	}

	if len(result.RouteTables) != 1 {
		return errors.New(fmt.Sprintf("found: %d vpc(s), should only find 1", len(result.RouteTables)))
	}

	r.InternetGateway.Id = result.RouteTables[0].RouteTableId
	return nil
}
