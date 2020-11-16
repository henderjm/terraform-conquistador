package resources

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
)

type client struct {
	envName      string
	internalTgw  string
	externalTgw  string
	outputFile   string
	awsResources AwsResources
	awsClient    *AWSClient
}

type AwsResources struct {
	networking networking
	publicRT   routeTable
	privateRT  routeTable
	natGateway natgw
	elb        elb
	elasticIp  eip
}

func NewClient(eN, i, e, oF string) client {
	fmt.Println("creating client")
	config := Config{
		AccessKey:     "access_key",
		SecretKey:     "secret_key",
		CredsFilename: "~/.aws/credentials",
		Region:        "eu-west-1",
	}
	aC, err := config.Client()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return client{
		envName:     eN,
		internalTgw: i,
		externalTgw: e,
		outputFile:  oF,
		awsClient:   aC,
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
	pubRT := NewRouteTable("Pub-RT")
	err = pubRT.Import(c)
	if err != nil {
		return err
	}
	privRT := NewRouteTable("PrivWeb-RT")
	err = privRT.Import(c)
	if err != nil {
		return err
	}

	// import loadbalancing resources
	lb := NewELB()
	err = lb.Import(c)
	if err != nil {
		return err
	}

	// import nat gateway
	ngw := NewNatGateway()
	err = ngw.Import(c)
	if err != nil {
		return err
	}
	eip := NewElasticIp()
	err = eip.Import(c)
	if err != nil {
		return err
	}

	// next resource
	c.awsResources = AwsResources{
		networking: n,
		publicRT:   pubRT,
		privateRT:  privRT,
		natGateway: ngw,
		elb:        lb,
		elasticIp:  eip,
	}
	return nil
}

func (c *client) UpdateTerraformStateFile() {
	// print networking
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' aws_vpc.base_vpc %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.networking.Vpc.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' aws_internet_gateway.ig %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.networking.Igw.Id))

	// print ALB
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' aws_lb.portal_lb %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.arn.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_subnet.lb_subnets[0]' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.subnets[0].Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_subnet.lb_subnets[1]' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elb.subnets[1].Id))

	// print Nat Gateway
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_nat_gateway.nat_gateway' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.natGateway.nat.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_eip.nat_eip' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.elasticIp.ip.Id))

	// print import route tables
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_route_table.public_subnet_rt' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.publicRT.table.Id))
	fmt.Printf("terraform import -state \"%s\" -var-file='./vars/%s-vars.tfvars' 'aws_route_table.private_subnet_rt' %s\n", c.outputFile, c.envName, aws.StringValue(c.awsResources.privateRT.table.Id))
}

func (c *client) CreateRouteTableVarFile() {
	f, err := os.OpenFile(fmt.Sprintf("%s-routetables.tfvars", c.envName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// private routing -  internal, external, nat
	f.WriteString(fmt.Sprintf("private_rt_internal_tgw_routes = {\"tgw\" = \"%s\", \"routes\" = %s}\n\n",
		c.internalTgw,
		printRoutes(c.awsResources.privateRT.InterTgwRoutes)))
	f.WriteString(fmt.Sprintf("private_rt_external_tgw_routes = {\"tgw\" = \"%s\", \"routes\" = %s}\n\n",
		c.externalTgw,
		printRoutes(c.awsResources.privateRT.ExternalTgwRoutes)))
	f.WriteString(fmt.Sprintf("private_rt_nat_routes = {\"nat\" = \"%s\", \"routes\" = %s}\n\n",
		aws.StringValue(c.awsResources.natGateway.nat.Id),
		printRoutes(c.awsResources.privateRT.NatRoutes)))
	// public routing - internal, external, igw
	f.WriteString(fmt.Sprintf("public_rt_internal_tgw_routes = {\"tgw\" = \"%s\", \"routes\" = %s}\n\n",
		c.internalTgw,
		printRoutes(c.awsResources.publicRT.InterTgwRoutes)))
	f.WriteString(fmt.Sprintf("public_rt_external_tgw_routes = {\"tgw\" = \"%s\", \"routes\" = %s}\n\n",
		c.externalTgw,
		printRoutes(c.awsResources.publicRT.ExternalTgwRoutes)))
	f.WriteString(fmt.Sprintf("public_rt_igw_routes = {\"igw\" = \"%s\", \"routes\" = %s}\n\n",
		aws.StringValue(c.awsResources.networking.Igw.Id),
		printRoutes(c.awsResources.publicRT.IgwRoutes)))

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func printRoutes(arr []route) string {
	var dests []string
	for _, r := range arr {
		dests = append(dests, fmt.Sprintf("\"%s\"", r.DestCidr))
	}

	return fmt.Sprintf("[%s]", strings.Join(dests, ","))
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
