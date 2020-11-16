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
	table             AWSResourceId
	routes            []*ec2.Route
	InterTgwRoutes    []route
	ExternalTgwRoutes []route
	NatRoutes         []route
	IgwRoutes         []route
	tag               string
}

type route struct {
	DestCidr  string
	GatewayID string
}

func NewRouteTable(t string) routeTable {
	return routeTable{
		tag: t,
	}
}

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
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("%s-%s", c.envName, r.tag))},
			},
		},
	}

	var result, err = conn.DescribeRouteTables(input)
	if err != nil {
		handleAWSError(err)
		return err
	}

	if len(result.RouteTables) != 1 {
		return errors.New(fmt.Sprintf("found: %d route table(s), should only find 1", len(result.RouteTables)))
	}

	r.table.Id = result.RouteTables[0].RouteTableId
	for _, awsRoute := range result.RouteTables[0].Routes {

		routeConf := route{
			DestCidr:  aws.StringValue(awsRoute.DestinationCidrBlock),
			GatewayID: "",
		}

		if awsRoute.TransitGatewayId != nil {
			routeConf.GatewayID = aws.StringValue(awsRoute.TransitGatewayId)
			if aws.StringValue(awsRoute.TransitGatewayId) == c.internalTgw {
				r.InterTgwRoutes = append(r.InterTgwRoutes, routeConf)
			} else if aws.StringValue(awsRoute.TransitGatewayId) == c.externalTgw {
				r.ExternalTgwRoutes = append(r.ExternalTgwRoutes, routeConf)
			} else {
				fmt.Printf(
					"error associating Transit Gateway `%s`. It matched neither internal: `%s` or external: `%s`\n",
					aws.StringValue(awsRoute.TransitGatewayId),
					c.internalTgw,
					c.externalTgw,
				)
			}
		}
		if awsRoute.NatGatewayId != nil {
			routeConf.GatewayID = aws.StringValue(awsRoute.NatGatewayId)
			r.NatRoutes = append(r.NatRoutes, routeConf)
		}
		if awsRoute.GatewayId != nil && aws.StringValue(awsRoute.GatewayId) != "local" {
			routeConf.GatewayID = aws.StringValue(awsRoute.GatewayId)
			r.IgwRoutes = append(r.IgwRoutes, routeConf)
		}
	}
	r.routes = result.RouteTables[0].Routes
	return nil
}
