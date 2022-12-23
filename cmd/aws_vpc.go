package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var vpcCMD = &cobra.Command{
	Use:   "vpc",
	Short: "Manually create / delete VPC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vpc called\n")
	},
}

var vpcCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create  VPC",
	Run: func(cmd *cobra.Command, args []string) {

		e.CreateByoipVpc()
		AddTag(e.SVC, e.VpcId, "Project", e.Project)

	},
}

var vpcDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete VPC",
	Run: func(cmd *cobra.Command, args []string) {

		e.DeleteVpc()

	},
}

var vpcDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe VPC",
	Run: func(cmd *cobra.Command, args []string) {

		e.DescribeVpc(vpcId)
	},
}

var vpcListCMD = &cobra.Command{
	Use:   "list",
	Short: "List VPC with current Project tags",
	Run: func(cmd *cobra.Command, args []string) {

		e.ListVpc()

	},
}

// Bundling functions

func (e *Environment) GetVpcId() error {
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	}

	result, err := e.SVC.DescribeVpcs(input)
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

	for _, vpcs := range result.Vpcs {
		//c.Environment[0].VpcId = *vpcs.VpcId
		e.VpcId = *vpcs.VpcId
	}

	if len(result.Vpcs) > 1 {
		err = xsIdErr
	} else if len(result.Vpcs) < 1 {
		err = noIdErr
	}

	return err
}

/*
func (e *Environment) DescribeVpc(vpcId string) error {

	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcId),
		},
	}

	result, err := e.SVC.DescribeVpcs(input)
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

	fmt.Printf("%v", result)
	return err

}
*/
func (e *Environment) CreateByoipVpc() error {

	input := &ec2.CreateVpcInput{
		CidrBlock:     aws.String("172.24.0.0/16"),
		Ipv6CidrBlock: aws.String(e.Byoip6cidr),
		Ipv6Pool:      aws.String(e.Byoip6pool),
		//Ipv6CidrBlock: aws.String("2a10:ba00:bee5:0000:0000:0000:0000:0000/56"),
		//Ipv6Pool:      aws.String("ipv6pool-ec2-0b40a8d4e1c7614d1"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("vpc"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("EnvTag"),
						Value: aws.String(e.EnvTag),
					},
				},
			},
		},
	}

	result, err := e.SVC.CreateVpc(input)
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

	if err == nil {
		//id = *result.Vpc.VpcId
		e.VpcId = *result.Vpc.VpcId
	}
	if verbose {
		fmt.Println(result)
	}

	return err

}

func (e *Environment) DescribeVpc(vpcId string) error {

	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcId),
		},
	}

	result, err := e.SVC.DescribeVpcs(input)
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

	fmt.Printf("%v", result)
	return err

}

func (e *Environment) ListVpc() error {

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	//fmt.Printf("%s", input)

	result, err := e.SVC.DescribeVpcs(input)
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

	//fmt.Printf("%v", result)
	fmt.Printf("%-30s %s\n", "VPC ID", "Tags")
	for _, vpcs := range result.Vpcs {
		vpcId := *vpcs.VpcId
		tags := ""
		if vpcs.Tags != nil {
			for _, tv := range result.Vpcs[0].Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
			}
		}
		fmt.Printf("%-30s %s\n", vpcId, tags)
	}

	return err

}

//func (e *Environment) DeleteVpc(vpcId string) error {
func (e *Environment) DeleteVpc() error {

	input := &ec2.DeleteVpcInput{
		//VpcId: aws.String(vpcId),
		VpcId: aws.String(e.VpcId),
	}

	result, err := e.SVC.DeleteVpc(input)
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
	return nil

}

func init() {

	//sess2, _ = CreateAwsSession()

	vpcCMD.AddCommand(vpcCreateCMD, vpcDeleteCMD, vpcDescribeCMD, vpcListCMD)
	vpcDeleteCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	vpcDescribeCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	vpcDeleteCMD.MarkFlagRequired("vpc-id")
	vpcDescribeCMD.MarkFlagRequired("vpc-id")
}
