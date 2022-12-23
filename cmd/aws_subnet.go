package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var subnetCMD = &cobra.Command{
	Use:   "subnet",
	Short: "Manually create / delete Subnet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("subnet called\n")

	},
}

var subnetCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create Subnet",
	Run: func(cmd *cobra.Command, args []string) {

		if e.VpcId == "" {
			e.VpcId = vpcId
		}
		e.CreateSubnet()
		AddTag(e.SVC, e.SubnetId, "Project", e.Project)

	},
}

var subnetDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Subnet",
	Run: func(cmd *cobra.Command, args []string) {

		e.DeleteSubnet()

	},
}

var subnetDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Subnet",
	Run: func(cmd *cobra.Command, args []string) {

		e.DescribeSubnet(subnetId)

	},
}

var subnetListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Subnets for Project",
	Run: func(cmd *cobra.Command, args []string) {

		e.ListSubnet()

	},
}

// Bundling functions

//func (e *Environment) GetSubnetId(svc *ec2.EC2) error {
func (e *Environment) GetSubnetId() error {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	}

	result, err := e.SVC.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	for _, subs := range result.Subnets {
		e.SubnetId = *subs.SubnetId
	}

	if len(result.Subnets) > 1 {
		err = xsIdErr
	} else if len(result.Subnets) < 1 {
		err = noIdErr
	}

	return err

}

func (e *Environment) CreateSubnet() error {

	// Make a /64 subnet from IPv6 Cidr
	var re = regexp.MustCompile(`/\d*$`)
	SubCidr := re.ReplaceAllString(e.Byoip6cidr, "/64")

	input := &ec2.CreateSubnetInput{
		//Ipv6CidrBlock: aws.String("2a10:ba00:bee5::/64"),
		Ipv6CidrBlock: aws.String(SubCidr),
		CidrBlock:     aws.String("172.24.0.0/24"),
		VpcId:         aws.String(e.VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("subnet"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("EnvTag"),
						Value: aws.String(e.EnvTag),
					},
				},
			},
		},
	}

	result, err := e.SVC.CreateSubnet(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	if err == nil {
		e.SubnetId = *result.Subnet.SubnetId
	}

	if verbose {
		fmt.Println(result)
	}
	return err

}

//func (e *Environment) DeleteSubnet(subnetId string) error {
func (e *Environment) DeleteSubnet() error {

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(e.SubnetId),
		//SubnetId: aws.String(subnetId),
	}

	result, err := e.SVC.DeleteSubnet(input)
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

	if verbose {
		fmt.Println(result)
	}
	return nil

}

func (e *Environment) ListSubnet() error {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	result, err := e.SVC.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	fmt.Printf("%-30s %-30s %-30s %-30s %s\n", "SUBNET-ID", "vpc-id", "CIDR", "CIDR6", "Tags")
	for _, subs := range result.Subnets {
		subnetId := *subs.SubnetId
		vpcId := *subs.VpcId
		cidr := ""
		cidr6 := ""
		if subs.CidrBlock != nil {
			cidr = *subs.CidrBlock
		}
		if subs.Ipv6CidrBlockAssociationSet != nil {
			for _, c6 := range subs.Ipv6CidrBlockAssociationSet {
				cidr6 += *c6.Ipv6CidrBlock + " "
			}
		}
		tags := ""
		if subs.Tags != nil {
			for _, tv := range subs.Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
			}
		}
		fmt.Printf("%-30s %-30s %-30s %-30s %s\n", subnetId, vpcId, cidr, cidr6, tags)
	}

	return err
}

func (e *Environment) DescribeSubnet(subnetId string) error {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			//aws.String(e.SubnetId),
			aws.String(subnetId),
		},
	}

	result, err := e.SVC.DescribeSubnets(input)
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

	fmt.Println(result)
	return err
}

func init() {
	subnetCMD.AddCommand(subnetCreateCMD, subnetDeleteCMD, subnetDescribeCMD, subnetListCMD)
	subnetCreateCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	subnetDeleteCMD.Flags().StringVarP(&subnetId, "subnet-id", "", "", "Subnet ID")
	subnetDescribeCMD.Flags().StringVarP(&subnetId, "subnet-id", "", "", "Subnet ID")
	subnetCreateCMD.MarkFlagRequired("vpc-id")
	subnetDeleteCMD.MarkFlagRequired("subnet-id")
	subnetDescribeCMD.MarkFlagRequired("subnet-id")
}
