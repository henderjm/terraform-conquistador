package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

type client struct {
	envName      string
	phase        string
	tags         []string
	outputFile   string
	awsResources AwsResources
	awsClient    *AWSClient
}

type AwsResources struct {
	networking networking
}

func NewClient(eN, p, t, oF string) *client {
	fmt.Println("creating client")
	tags := strings.Split(t, ",")
	config := Config{
		AccessKey:     "access_key",
		SecretKey:     "secret_key",
		CredsFilename: "~/.aws/credentials",
		Region:        "eu-west-1",
	}
	aC, err := config.Client()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return &client{
		envName:    eN,
		phase:      p,
		tags:       tags,
		outputFile: oF,
		awsClient:  aC,
	}
}

func (c *client) ImportTerraformResources() error {
	var err error
	// import networking
	n := NewNetworking()
	c.awsResources.networking, err = n.Import(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateTerraformStateFile() {
	fmt.Printf("terraform import -var-file=vars.tfvars aws_vpc.base_vpc %s\n", aws.StringValue(c.awsResources.networking.Vpc.Id))
	fmt.Printf("terraform import -var-file=vars.tfvars aws_internet_gateway.ig %s\n", aws.StringValue(c.awsResources.networking.Ig.Id))
}
