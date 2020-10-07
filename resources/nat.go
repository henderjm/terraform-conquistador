package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
)

/*
Nat Gateway
*/

type natgw struct {
	nat    AWSResourceId
	subnet AWSResourceId
}

func NewNatGateway() natgw { return natgw{} }

func (n *natgw) Import(c *client) error {
	err := n.importNat(c)
	if err != nil {
		return err
	}
	return nil
}

func (n *natgw) importNat(c *client) error {
	fmt.Println("searching-for-nat-gateway")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{ // TODO: PR To fix this
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
		},
	}

	result, err := conn.DescribeNatGateways(input)
	if err != nil {
		handleAWSError(err)
		return err
	}

	if len(result.NatGateways) != 1 {
		return errors.New(fmt.Sprintf("found: %d ig(s), should only find 1", len(result.NatGateways)))
	}

	ngw := result.NatGateways[0]

	n.nat.Id = ngw.NatGatewayId
	n.subnet.Id = ngw.SubnetId

	return nil
}
