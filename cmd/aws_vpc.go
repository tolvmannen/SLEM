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

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_, vpcId := CreateVpc(svc)
		fmt.Printf("VPC created: %s\n", vpcId)

	},
}

var vpcDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete VPC",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteVpc(svc, vpcId)

	},
}

var vpcDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe VPC",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeVpc(svc, vpcId)

	},
}

var vpcListCMD = &cobra.Command{
	Use:   "list",
	Short: "List VPC with current Project tags",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		_ = ListVpc(svc)

	},
}

func CreateVpc(svc *ec2.EC2) (error, string) {

	input := &ec2.CreateVpcInput{
		CidrBlock:     aws.String("172.24.0.0/16"),
		Ipv6CidrBlock: aws.String("2a10:ba00:bee5:0000:0000:0000:0000:0000/56"),
		Ipv6Pool:      aws.String("ipv6pool-ec2-0b40a8d4e1c7614d1"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("vpc"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	}

	result, err := svc.CreateVpc(input)
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
	if result != nil {
		id = *result.Vpc.VpcId
	}
	if verbose {
		fmt.Println(result)
	}
	return err, id

}

func DeleteVpc(svc *ec2.EC2, vpcId string) error {

	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcId),
	}

	result, err := svc.DeleteVpc(input)
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

func IsValidVpcId(svc *ec2.EC2, vpcId string) bool {
	result := GetVpcId(svc, vpcId)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

func GetVpcId(svc *ec2.EC2, vpcId string) string {

	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcId),
		},
	}

	result, err := svc.DescribeVpcs(input)
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

	//fmt.Printf("\nid Exists: %v\n", result.Vpcs[0].VpcId)

	var id string
	if len(result.Vpcs) > 0 {
		id = *result.Vpcs[0].VpcId
	}

	return id

}

func DescribeVpc(svc *ec2.EC2, vpcId string) error {

	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcId),
		},
	}

	result, err := svc.DescribeVpcs(input)
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

func ListVpc(svc *ec2.EC2) error {

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	//fmt.Printf("%s", input)

	result, err := svc.DescribeVpcs(input)
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

func init() {
	vpcCMD.AddCommand(vpcCreateCMD, vpcDeleteCMD, vpcDescribeCMD, vpcListCMD)
	vpcDeleteCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	vpcDescribeCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	vpcDeleteCMD.MarkFlagRequired("vpc-id")
	vpcDescribeCMD.MarkFlagRequired("vpc-id")
}
