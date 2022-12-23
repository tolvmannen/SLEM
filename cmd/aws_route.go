package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var routeCMD = &cobra.Command{
	Use:   "route",
	Short: "Manually create / delete Route (Add to/Remove from to Route Table)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("route-table called\n")

	},
}

var routeCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create Route",
	Run: func(cmd *cobra.Command, args []string) {

		if e.IgwId == "" {
			e.IgwId = igwId
		}
		if e.RtbId == "" {
			e.RtbId = rtbId
		}
		e.CreateRoute()
		//if err == nil {
		//	fmt.Printf("Route added to Route Table: %s\n", rtbId)
		//}

	},
}

var routeDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Route",
	Run: func(cmd *cobra.Command, args []string) {

		e.DeleteRoute(rtbId)

	},
}

// Bundling functions

func (e *Environment) CreateRoute() error {

	input := &ec2.CreateRouteInput{
		GatewayId:            aws.String(e.IgwId),
		RouteTableId:         aws.String(e.RtbId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		//DestinationIpv6CidrBlock: aws.String("::0/0"),
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

func (e *Environment) DeleteRoute(rtbId string) error {

	input := &ec2.DeleteRouteInput{
		RouteTableId:         aws.String(rtbId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		//DestinationIpv6CidrBlock: aws.String("::0/0"),
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
	return err

}

func init() {
	routeCMD.AddCommand(routeCreateCMD, routeDeleteCMD)
	routeCreateCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	routeCreateCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	routeDeleteCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")

	routeCreateCMD.MarkFlagRequired("rtb-id")
	routeCreateCMD.MarkFlagRequired("igw-id")
	routeDeleteCMD.MarkFlagRequired("rtb-id")
}
