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

		e.CreateIGW()
		AddTag(e.SVC, e.IgwId, "Project", e.Project)

	},
}

var igwDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		e.DeleteIGW()

	},
}

var igwAttachCMD = &cobra.Command{
	Use:   "attach",
	Short: "Manually attach Internet Gateway to VPV",
	Run: func(cmd *cobra.Command, args []string) {

		e.AttachIGW()

	},
}

var igwDetatchCMD = &cobra.Command{
	Use:   "detach",
	Short: "Manually detatch Internet Gateway from VPC",
	Run: func(cmd *cobra.Command, args []string) {

		if e.IgwId == "" {
			e.IgwId = igwId
		}
		if e.VpcId == "" {
			e.VpcId = vpcId
		}
		e.DetatchIGW()

	},
}

var igwDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		e.DescribeIGW(igwId)

	},
}

var igwListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Internet Gateway for project",
	Run: func(cmd *cobra.Command, args []string) {
		e.ListIGW()

	},
}

// Bundling function§
//func (c *MainConfig) GetIgwId() error {
func (e *Environment) GetIgwId() error {

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	}

	result, err := e.SVC.DescribeInternetGateways(input)
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

	for _, igws := range result.InternetGateways {
		e.IgwId = *igws.InternetGatewayId
	}

	if len(result.InternetGateways) > 1 {
		err = xsIdErr
	} else if len(result.InternetGateways) < 1 {
		err = noIdErr
	}

	return err
}

func (e *Environment) CreateIGW() error {

	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("internet-gateway"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("EnvTag"),
						Value: aws.String(e.EnvTag),
					},
				},
			},
		},
	}

	result, err := e.SVC.CreateInternetGateway(input)
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

	if result != nil {
		e.IgwId = *result.InternetGateway.InternetGatewayId
	}
	if verbose {
		fmt.Println(result)
	}
	return err

}

//func (c *MainConfig) DeleteIGW(igwId string) error {
func (e *Environment) DeleteIGW() error {

	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(e.IgwId),
	}

	// Success returns nil for result and err
	result, err := e.SVC.DeleteInternetGateway(input)
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

func (e *Environment) AttachIGW() error {

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(e.IgwId),
		VpcId:             aws.String(e.VpcId),
	}

	// Success returns nil for result and err
	result, err := e.SVC.AttachInternetGateway(input)
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

func (e *Environment) DetatchIGW() error {

	input := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(e.IgwId),
		VpcId:             aws.String(e.VpcId),
	}

	// Success returns nil for result and err
	result, err := e.SVC.DetachInternetGateway(input)
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

func (e *Environment) ListIGW() error {

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	result, err := e.SVC.DescribeInternetGateways(input)
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

func (e *Environment) DescribeIGW(igwId string) error {

	input := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{
			//aws.String(e.IgwId),
			aws.String(igwId),
		},
	}

	result, err := e.SVC.DescribeInternetGateways(input)
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
