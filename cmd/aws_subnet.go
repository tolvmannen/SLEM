package cmd

import (
	"fmt"

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

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		err, subnetId := CreateSubnet(svc, vpcId)
		if err != nil {
			fmt.Printf("Failed to create Subnet, %v", err)
		} else {
			fmt.Printf("Subnet created: %s\n", subnetId)
		}

	},
}

var subnetDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Subnet",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteSubnet(svc, subnetId)

	},
}

var subnetDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Subnet",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeSubnet(svc, subnetId)

	},
}

var subnetListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Subnets for Project",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = ListSubnet(svc)

	},
}

func CreateSubnet(svc *ec2.EC2, vpcId string) (error, string) {

	input := &ec2.CreateSubnetInput{
		Ipv6CidrBlock: aws.String("2a10:ba00:bee5:0000::/64"),
		CidrBlock:     aws.String("172.24.0.0/24"),
		VpcId:         aws.String(vpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("subnet"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	}

	result, err := svc.CreateSubnet(input)
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

	subnetId := *result.Subnet.SubnetId

	//fmt.Println(result)
	return err, subnetId

}

func DeleteSubnet(svc *ec2.EC2, subnetId string) error {

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetId),
	}

	result, err := svc.DeleteSubnet(input)
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

func DescribeSubnet(svc *ec2.EC2, subnetId string) error {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(subnetId),
		},
	}

	result, err := svc.DescribeSubnets(input)
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

func ListSubnet(svc *ec2.EC2) error {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	result, err := svc.DescribeSubnets(input)
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

func GetSubnetId(svc *ec2.EC2, subnetId string) string {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(subnetId),
		},
	}

	result, err := svc.DescribeSubnets(input)
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

	var id string
	if len(result.Subnets) > 0 {
		id = *result.Subnets[0].SubnetId
	}

	//fmt.Println(result)
	return id
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
