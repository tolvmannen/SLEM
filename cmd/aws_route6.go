package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var route6CMD = &cobra.Command{
	Use:   "route6",
	Short: "Manually create / delete IPv6 Route (Add to/Remove from to Route Table)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("route-table called\n")

	},
}

var route6CreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create IPv6 Route",
	Run: func(cmd *cobra.Command, args []string) {

		if e.IgwId == "" {
			e.IgwId = igwId
		}
		if e.RtbId == "" {
			e.RtbId = rtbId
		}
		e.CreateRoute6()
		//if err == nil {
		//	fmt.Printf("Route added to Route Table: %s\n", rtbId)
		//}

	},
}

var route6DeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete IPv6 Route",
	Run: func(cmd *cobra.Command, args []string) {

		e.DeleteRoute6(rtbId)

	},
}

// Bundling Ffunctions

func (e *Environment) CreateRoute6() error {

	input := &ec2.CreateRouteInput{
		GatewayId:    aws.String(e.IgwId),
		RouteTableId: aws.String(e.RtbId),
		//DestinationCidrBlock: aws.String("0.0.0.0/0"),
		DestinationIpv6CidrBlock: aws.String("::/0"),
	}

	result, err := e.SVC.CreateRoute(input)
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
	return err

}

func (e *Environment) DeleteRoute6(rtbId string) error {

	input := &ec2.DeleteRouteInput{
		RouteTableId: aws.String(rtbId),
		//DestinationCidrBlock: aws.String("0.0.0.0/0"),
		DestinationIpv6CidrBlock: aws.String("::/0"),
	}

	result, err := e.SVC.DeleteRoute(input)
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

func init() {
	route6CMD.AddCommand(route6CreateCMD, route6DeleteCMD)
	route6CreateCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	route6CreateCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	route6DeleteCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")

	route6CreateCMD.MarkFlagRequired("rtb-id")
	route6CreateCMD.MarkFlagRequired("igw-id")
	route6DeleteCMD.MarkFlagRequired("rtb-id")
}
