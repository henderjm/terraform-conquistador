package main

import (
	"fmt"
	"os"

	"github.com/henderjm/terraform-conquistador/resources"
)

func main() {
	envName := os.Args[1]
	phase := os.Args[2]
	tags := os.Args[3]
	stateFile := os.Args[4]

	c := resources.NewClient(envName, phase, tags, stateFile)
	err := c.ImportTerraformResources() // Change client function to run
	// Change VPC to Network which contains importing VPC, Subnets, Security Groups
	// Add EC2 to import EC2, LBs, Target Groups,

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("...updating terraform state file: %s\n", stateFile)
	c.UpdateTerraformStateFile()
}
