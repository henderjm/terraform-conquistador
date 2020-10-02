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

type elb struct {
	alb AWSResourceId
}

func NewELB() elb { return elb{} }

func (e *elb) Import(c *client) (elb, error) {
	a, err := importALB(c)
	if err != nil {
		return elb{}, err
	}
	return elb{
		alb: a,
	}, nil
}

func importALB(c *client) (AWSResourceId, error) {
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
		return AWSResourceId{}, err
	}

	if len(result.LoadBalancers) != 1 {
		return AWSResourceId{}, errors.New(fmt.Sprintf("found: %d ig(s), should only find 1", len(result.LoadBalancers)))
	}

	alb := AWSResourceId{
		Id: result.LoadBalancers[0].LoadBalancerArn,
	}

	return alb, nil
}
