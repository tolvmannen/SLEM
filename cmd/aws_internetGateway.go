package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var igwCMD = &cobra.Command{
	Use:   "igw",
	Short: "Manually create / delete VPC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("igw called\n")
	},
}

var igwCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		err, igwId := CreateInternetGateway(svc)
		if err == nil {
			fmt.Printf("Internet Gateway created: %s\n", igwId)
		}

	},
}

var igwDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		err = DeleteInternetGateway(svc, igwId)
		if err == nil {
			fmt.Printf("Internet Gateway deleted (%s)\n", igwId)
		}

	},
}

var igwAttachCMD = &cobra.Command{
	Use:   "attach",
	Short: "Manually attach Internet Gateway to VPV",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		err = AttachInternetGateway(svc, igwId, vpcId)
		if err == nil {
			fmt.Printf("Internet Gateway %s attached to %s\n", igwId, vpcId)
		}

	},
}

var igwDetatchCMD = &cobra.Command{
	Use:   "detach",
	Short: "Manually detatch Internet Gateway from VPC",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		err = DetatchInternetGateway(svc, igwId, vpcId)
		if err == nil {
			fmt.Printf("Internet Gateway %s detached from %s\n", igwId, vpcId)
		}

	},
}

var igwDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeInternetGateway(svc, igwId)

	},
}

var igwListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Internet Gateway for project",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = ListInternetGateway(svc)

	},
}

func CreateInternetGateway(svc *ec2.EC2) (error, string) {

	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("internet-gateway"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	}

	result, err := svc.CreateInternetGateway(input)
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
		id = *result.InternetGateway.InternetGatewayId
	}

	//fmt.Println(result)
	return err, id

}

func DeleteInternetGateway(svc *ec2.EC2, igwId string) error {

	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
	}

	// Success returns nil for result and err
	result, err := svc.DeleteInternetGateway(input)
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
		// result empty, but print anyway
		fmt.Println(result)
	}
	return err

}

func AttachInternetGateway(svc *ec2.EC2, igwId, vpcId string) error {

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	}

	// Success returns nil for result and err
	result, err := svc.AttachInternetGateway(input)
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

	// result empty, but print anyway (for now)
	if verbose {
		fmt.Println(result)
	}
	return err

}

func DetatchInternetGateway(svc *ec2.EC2, igwId, vpcId string) error {

	input := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	}

	// Success returns nil for result and err
	result, err := svc.DetachInternetGateway(input)
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
		// result empty, but print anyway
		fmt.Println(result)
	}
	return err

}

func ListInternetGateway(svc *ec2.EC2) error {

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	result, err := svc.DescribeInternetGateways(input)
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

	fmt.Printf("%-30s %-30s %s\n", "IGW-ID", "Attachments", "Tags")
	for _, igws := range result.InternetGateways {
		igwId := *igws.InternetGatewayId
		vpcId := ""
		if igws.Attachments != nil {
			vpcId = *igws.Attachments[0].VpcId
			if len(igws.Attachments) > 1 {
				vpcId += "**"
			}
		}
		tags := ""
		if igws.Tags != nil {
			for _, tv := range igws.Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
			}
		}
		fmt.Printf("%-30s %-30s %s\n", igwId, vpcId, tags)

	}

	return err
}

func DescribeInternetGateway(svc *ec2.EC2, igwId string) error {

	input := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{
			aws.String(igwId),
		},
	}

	result, err := svc.DescribeInternetGateways(input)
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

func init() {
	igwCMD.AddCommand(igwCreateCMD, igwDeleteCMD, igwAttachCMD, igwDetatchCMD, igwDescribeCMD, igwListCMD)
	igwDeleteCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	igwAttachCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	igwAttachCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	igwDetatchCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	igwDetatchCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	igwDescribeCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	igwAttachCMD.MarkFlagRequired("igw-id")
	igwAttachCMD.MarkFlagRequired("vpc-id")
	igwDetatchCMD.MarkFlagRequired("igw-id")
	igwDetatchCMD.MarkFlagRequired("vpc-id")
	igwDescribeCMD.MarkFlagRequired("igw-id")

}
