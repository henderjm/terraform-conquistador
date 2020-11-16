package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
)

/*
eip
*/

type eip struct {
	ip AWSResourceId
}

func NewElasticIp() eip { return eip{} }

func (e *eip) Import(c *client) error {
	err := e.importEIP(c)
	if err != nil {
		return err
	}
	return nil
}

func (e *eip) importEIP(c *client) error {
	fmt.Println("searching-for-nat-gateway")
	conn := c.awsClient.ec2conn
	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{ // TODO: PR To fix this
			{
				Name:   aws.String("tag:VPC"),
				Values: []*string{aws.String(fmt.Sprintf("%s", c.envName))},
			},
		},
	}

	result, err := conn.DescribeAddresses(input)
	if err != nil {
		handleAWSError(err)
		return err
	}

	if len(result.Addresses) != 1 {
		return errors.New(fmt.Sprintf("found: %d elastic ip(s), should only find 1", len(result.Addresses)))
	}

	e.ip.Id = result.Addresses[0].PublicIp

	return nil
}
