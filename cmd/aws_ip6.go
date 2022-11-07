package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var eip6CMD = &cobra.Command{
	Use:   "eip6",
	Short: "Manually manage Elastic IP Addresses (for BYOIP)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ip6 called\n")

	},
}

var eipAssignIpv6CMD = &cobra.Command{
	Use:   "assign",
	Short: "Associate an EIP with an EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = AssignIpv6(svc, eniId, ip6)
		fmt.Printf("Associated %s", eipassocId)

	},
}

func AssignIpv6(svc *ec2.EC2, eniId, ip6 string) string {

	input := &ec2.AssignIpv6AddressesInput{
		NetworkInterfaceId: aws.String(eniId),
		Ipv6Addresses: []*string{
			aws.String("2a10:ba00:bee5::4"),
		},
	}

	result, err := svc.AssignIpv6Addresses(input)
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
	if err == nil {
		id = *result.AssignedIpv6Addresses[0]
	}
	//fmt.Printf("%T \n %#v", result, result)
	//fmt.Println(result)
	return id
}

func init() {
	eip6CMD.AddCommand(eipListCMD, eipAssignIpv6CMD)

	eipAssignIpv6CMD.Flags().StringVarP(&ip6, "ip6", "", "", "IPv6 address")
	eipAssignIpv6CMD.Flags().StringVarP(&eniId, "eni-id", "", "", "EC2 Network Interface ID")

	eipAssignIpv6CMD.MarkFlagRequired("ip6")
	eipAssignIpv6CMD.MarkFlagRequired("eni-id")
}
