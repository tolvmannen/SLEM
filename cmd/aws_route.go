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

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		err = CreateRoute(svc, igwId, rtbId)

		if err == nil {
			fmt.Printf("Route added to Route Table: %s\n", rtbId)
		}

	},
}

var routeCreate6CMD = &cobra.Command{
	Use:   "create6",
	Short: "Manually create Route",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		err = CreateRoute(svc, igwId, rtbId)

		if err != nil {
			fmt.Printf("Failed to create Toute Table, %v", err)
		} else {
			fmt.Printf("Added Rout to Route Table: %s\n", rtbId)
		}

	},
}

var routeDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Route",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteRoute(svc, rtbId)

	},
}

var routeDelete6CMD = &cobra.Command{
	Use:   "delete6",
	Short: "Manually delete Route",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteRoute6(svc, rtbId)

	},
}

func CreateRoute(svc *ec2.EC2, igwId, rtbId string) error {

	input := &ec2.CreateRouteInput{
		GatewayId:            aws.String(igwId),
		RouteTableId:         aws.String(rtbId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		//DestinationIpv6CidrBlock: aws.String("::0/0"),
	}

	result, err := svc.CreateRoute(input)
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

func DeleteRoute(svc *ec2.EC2, rtbId string) error {

	input := &ec2.DeleteRouteInput{
		RouteTableId:         aws.String(rtbId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		//DestinationIpv6CidrBlock: aws.String("::0/0"),
	}

	result, err := svc.DeleteRoute(input)
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
	routeCMD.AddCommand(routeCreateCMD, routeDeleteCMD)
	routeCreateCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")
	routeCreateCMD.Flags().StringVarP(&igwId, "igw-id", "", "", "Internet Gateway ID")
	routeDeleteCMD.Flags().StringVarP(&rtbId, "rtb-id", "", "", "Routing Table ID")

	routeCreateCMD.MarkFlagRequired("rtb-id")
	routeCreateCMD.MarkFlagRequired("igw-id")
	routeDeleteCMD.MarkFlagRequired("rtb-id")
}
