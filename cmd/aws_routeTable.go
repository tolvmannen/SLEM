package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var rtbCMD = &cobra.Command{
	Use:   "route-table",
	Short: "Manually create / delete Route Table",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("route-table called\n")

	},
}

var rtbCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create Route Table",
	Run: func(cmd *cobra.Command, args []string) {

		//cf.Environment.CreateRTB(cf.SVC)
		e.CreateRTB()
		AddTag(e.SVC, e.RtbId, "Project", e.Project)

	},
}

var rtbDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Route Table",
	Run: func(cmd *cobra.Command, args []string) {

		if e.RtbId == "" {
			e.RtbId = rtbId
		}
		e.DeleteRTB()
	},
}

var rtbAssociateCMD = &cobra.Command{
	Use:   "associate",
	Short: "Associate a Route Table to an Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		if e.RtbId == "" {
			e.RtbId = rtbId
		}
		if e.SubnetId == "" {
			e.SubnetId = subnetId
		}
		e.AssociateRTB()

	},
}

var rtbDisassociateCMD = &cobra.Command{
	Use:   "disassociate",
	Short: "Disassociate a Route Table from an Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		if e.RtbassocId == "" {
			e.RtbassocId = rtbassocId
		}
		e.DisassociateRTB()

	},
}

var rtbListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Route Tables for Project",
	Run: func(cmd *cobra.Command, args []string) {

		e.ListRTB()

	},
}

var rtbDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe a Route Table",
	Run: func(cmd *cobra.Command, args []string) {

		e.DescribeRTB(rtbId)

	},
}

// Bundling Functions

func (e *Environment) GetRtbId() error {

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	}

	result, err := e.SVC.DescribeRouteTables(input)
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

	for _, rtbs := range result.RouteTables {
		e.RtbId = *rtbs.RouteTableId
	}

	if len(result.RouteTables) > 1 {
		err = xsIdErr
	} else if len(result.RouteTables) < 1 {
		err = noIdErr
	}

	return err

}

func (e *Environment) CreateRTB() error {

	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(e.VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("route-table"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("EnvTag"),
						Value: aws.String(e.EnvTag),
					},
				},
			},
		},
	}

	result, err := e.SVC.CreateRouteTable(input)
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
		e.RtbId = *result.RouteTable.RouteTableId
	}

	return err

}

func (e *Environment) GetRtbassocId() error {

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(e.RtbId),
		},
	}

	result, err := e.SVC.DescribeRouteTables(input)
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

	for _, rtbs := range result.RouteTables {
		if rtbs.Associations != nil {
			e.RtbassocId = *rtbs.Associations[0].RouteTableAssociationId
			if len(rtbs.Associations) > 1 {
				err = xsIdErr
			}
		} else {
			err = noIdErr
		}
	}

	return err

}

func (e *Environment) AssociateRTB() error {

	input := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(e.RtbId),
		SubnetId:     aws.String(e.SubnetId),
	}

	result, err := e.SVC.AssociateRouteTable(input)
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
		e.RtbassocId = *result.AssociationId
	}
	if verbose {
		fmt.Println(result)
	}
	return err

}

func (e *Environment) DisassociateRTB() error {

	input := &ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(e.RtbassocId),
	}

	result, err := e.SVC.DisassociateRouteTable(input)
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

	if verbose {
		fmt.Println(result)
	}
	return nil

}

func (e *Environment) DeleteRTB() error {

	input := &ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(e.RtbId),
	}

	result, err := e.SVC.DeleteRouteTable(input)
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

	if verbose {
		fmt.Println(result)
	}
	return err

}

func (e *Environment) ListRTB() error {

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	result, err := e.SVC.DescribeRouteTables(input)
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

	fmt.Printf("%-30s %-30s %-30s %-30s %s\n", "RTB-ID", "vpc-id", "subnet-id", "rtbassoc-id", "Tags")
	for _, rtbs := range result.RouteTables {
		rtbId := *rtbs.RouteTableId
		vpcId := *rtbs.VpcId
		subnetId := ""
		rtbassocId := ""
		if rtbs.Associations != nil {
			subnetId = *rtbs.Associations[0].SubnetId
			rtbassocId = *rtbs.Associations[0].RouteTableAssociationId
			// Note that just the first entry is shown if there are several..
			// Maybe there never are?
			if len(rtbs.Associations) > 1 {
				rtbassocId += "**"
				subnetId += "**"
			}
		}
		tags := ""
		if rtbs.Tags != nil {
			for _, tv := range rtbs.Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
			}
		}
		fmt.Printf("%-30s %-30s %-30s %-30s %s\n", rtbId, vpcId, subnetId, rtbassocId, tags)

	}

	return err
}

func (e *Environment) DescribeRTB(rtbId string) error {

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(rtbId),
		},
	}

	result, err := e.SVC.DescribeRouteTables(input)
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
	rtbCMD.AddCommand(rtbCreateCMD, rtbDeleteCMD, rtbAssociateCMD, rtbDisassociateCMD, rtbDescribeCMD, rtbListCMD)
	rtbCreateCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	rtbDeleteCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	rtbDescribeCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	rtbAssociateCMD.Flags().StringVarP(&subnetId, "subnet-id", "", "", "Subnet ID")
	rtbAssociateCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	rtbDisassociateCMD.Flags().StringVarP(&rtbassocId, "rtbassoc-id", "", "", "Routing Table Association ID")
	rtbCreateCMD.MarkFlagRequired("vpc-id")
	rtbDeleteCMD.MarkFlagRequired("rtb-id")
	rtbDescribeCMD.MarkFlagRequired("rtb-id")
	rtbAssociateCMD.MarkFlagRequired("subnet-id")
	rtbAssociateCMD.MarkFlagRequired("rtb-id")
	rtbDisassociateCMD.MarkFlagRequired("rtbassoc-id")
}
