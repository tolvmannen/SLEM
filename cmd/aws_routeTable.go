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

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		err, rtbId := CreateRT(svc, vpcId)
		if err != nil {
			fmt.Printf("Failed to create Toute Table, %v", err)
		} else {
			fmt.Printf("RT created: %s\n", rtbId)
		}

	},
}

var rtbDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Route Table",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteRT(svc, rtbId)

	},
}

var rtbAssociateCMD = &cobra.Command{
	Use:   "associate",
	Short: "Associate a Route Table to an Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		err, rtbassocId = AssociateRT(svc, rtbId, subnetId)
		if err == nil {
			fmt.Printf("Route Table %s associated with Subnet %s : (%s)\n", rtbId, subnetId, rtbassocId)
		}

	},
}

var rtbDisassociateCMD = &cobra.Command{
	Use:   "disassociate",
	Short: "Disassociate a Route Table from an Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DisassociateRT(svc, rtbassocId)

	},
}

var rtbListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Route Tables for Project",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = ListRTB(svc)

	},
}

var rtbDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe a Route Table",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeRT(svc, rtbId)

	},
}

func CreateRT(svc *ec2.EC2, vpcId string) (error, string) {

	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("route-table"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	}

	result, err := svc.CreateRouteTable(input)
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

	var rtbId string
	if result != nil {
		rtbId = *result.RouteTable.RouteTableId
	}

	return err, rtbId

}

func DeleteRT(svc *ec2.EC2, rtbId string) error {

	input := &ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(rtbId),
	}

	result, err := svc.DeleteRouteTable(input)
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

func AssociateRT(svc *ec2.EC2, rtbId, subnetId string) (error, string) {

	input := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rtbId),
		SubnetId:     aws.String(subnetId),
	}

	result, err := svc.AssociateRouteTable(input)
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

	var rtbassocId string
	if result != nil {
		rtbassocId = *result.AssociationId
	}
	//fmt.Println(result)
	return err, rtbassocId

}

func DisassociateRT(svc *ec2.EC2, rtbassocId string) error {

	input := &ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(rtbassocId),
	}

	result, err := svc.DisassociateRouteTable(input)
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

func DescribeRT(svc *ec2.EC2, rtbId string) error {

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(rtbId),
		},
	}

	result, err := svc.DescribeRouteTables(input)
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

func ListRTB(svc *ec2.EC2) error {

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	result, err := svc.DescribeRouteTables(input)
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

func GetRtbId(svc *ec2.EC2, rtbId string) string {

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(rtbId),
		},
	}

	result, err := svc.DescribeRouteTables(input)
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
	if len(result.RouteTables) > 0 {
		id = *result.RouteTables[0].RouteTableId
	}

	fmt.Println(result)
	return id
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
