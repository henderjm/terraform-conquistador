package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/aws/aws-sdk-go/aws"
)

/*
ALB
*/
type lb struct {
	arn     AWSResourceId
	subnets []AWSResourceId
}

type elb struct {
	lb
}

func NewELB() elb { return elb{} }

func (e *elb) Import(c *client) error {
	err := e.importALB(c)
	if err != nil {
		return err
	}
	return nil
}

func (e *elb) importALB(c *client) error {
	fmt.Println("searching-for-alb")
	conn := c.awsClient.elbv2conn
	input := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(fmt.Sprintf("%s-Public-ALB", c.envName)),
		},
	}

	result, err := conn.DescribeLoadBalancers(input)
	if err != nil {
		handleAWSError(err)
		return err
	}

	if len(result.LoadBalancers) != 1 {
		return errors.New(fmt.Sprintf("found: %d ig(s), should only find 1", len(result.LoadBalancers)))
	}

	var subnets []AWSResourceId
	for _, s := range result.LoadBalancers[0].AvailabilityZones {
		r := AWSResourceId{
			Id: s.SubnetId,
		}
		subnets = append(subnets, r)
	}

	id := AWSResourceId{
		Id: result.LoadBalancers[0].LoadBalancerArn,
	}

	e.arn = id
	e.subnets = subnets

	return nil
}
