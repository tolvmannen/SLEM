package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var eipCMD = &cobra.Command{
	Use:   "eip",
	Short: "Manually manage Elastic IP Addresses (for BYOIP)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("eip called\n")

	},
}

var eipListCMD = &cobra.Command{
	Use:   "list",
	Short: "List EIP addresses (with usage status)",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		if ip4 != "" {
			_ = ListAddressesIp4(svc, ip4)
		} else {
			_ = ListAddresses(svc)
		}

	},
}

var eipAllocateCMD = &cobra.Command{
	Use:   "allocate",
	Short: "Create an IP-allocaion",
	Run: func(cmd *cobra.Command, args []string) {

		if !ValidIPAddress(ip4) {
			exitErrorf("%s is not a valid IPv4 address\n", ip4)
		}

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		svc := ec2.New(sess) //  type *ec2.EC2

		var eipallocId string
		err, eipallocId = AllocateAddress(svc, ip4)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%-20s Allocated (%s)\n", ip4, eipallocId)

		}

	},
}

var eipReleaseCMD = &cobra.Command{
	Use:   "release",
	Short: "Remove an IP-allocaion",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}
		svc := ec2.New(sess)

		err = ReleaseAddress(svc, eipallocId)
		if err == nil {
			fmt.Printf("%s released\n", eipallocId)
		}

	},
}

var eipAssociateCMD = &cobra.Command{
	Use:   "associate",
	Short: "Associate an EIP with an EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		eipassocId := AssociateAddress(svc, instanceId, eipallocId)
		fmt.Printf("IP Associated: (%s)\n", eipassocId)

	},
}

var eipDisassociateCMD = &cobra.Command{
	Use:   "disassociate",
	Short: "Disasssociate an EIP from an EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DisassociateAddress(svc, eipassocId)

	},
}

var eipDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe EIP Address",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeAddress(svc, ip4)

	},
}

func DescribeAddress(svc *ec2.EC2, ip4 string) error {
	// Make the API request to EC2 filtering for the addresses in the
	// account's VPC.
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("public-ip"),
				Values: aws.StringSlice([]string{ip4}),
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to get elastic IP address, %v", err)
	}

	fmt.Printf("\n%v\n", result)
	//ret := result.Addresses
	return err
}

func ListAddresses(svc *ec2.EC2) error {
	// Make the API request to EC2 filtering for the addresses in the
	// account's VPC.
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to get elastic IP address, %v", err)
	}

	fmt.Printf("%-20s %-30s %-30s %s\n", "PublicIp", "AllocationId", "AssociationId", "Tags")
	if len(result.Addresses) > 0 {
		for _, v := range result.Addresses {
			PrintAddrInfo(v)
		}
	}

	return err
}

func ListAddressesIp4(svc *ec2.EC2, ip4 string) error {
	// Make the API request to EC2 filtering for the addresses in the
	// account's VPC.
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("public-ip"),
				Values: aws.StringSlice([]string{ip4}),
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to get elastic IP address, %v", err)
	}

	fmt.Printf("%-20s %-30s %-30s %s\n", "PublicIp", "AllocationId", "AssociationId", "Tags")
	if len(result.Addresses) > 0 {
		for _, v := range result.Addresses {
			PrintAddrInfo(v)
		}
	}

	return err
}

/*
func DescribeAddress(svc *ec2.EC2, ip4 string) []*ec2.Address {
	// Make the API request to EC2 filtering for the addresses in the
	// account's VPC.
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("public-ip"),
				Values: aws.StringSlice([]string{ip4}),
				//Name:   aws.String("domain"),
				//Values: aws.StringSlice([]string{"vpc"}),
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to get elastic IP address, %v", err)
	}

	//fmt.Printf("\n%+v\n", result)
	ret := result.Addresses
	return ret
}
*/

func AllocateAddress(svc *ec2.EC2, ip4 string) (error, string) {
	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain:  aws.String("vpc"),
		Address: aws.String(ip4),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("elastic-ip"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	})
	var eipallocId string
	if allocRes.AllocationId != nil {
		eipallocId = *allocRes.AllocationId
	}
	return err, eipallocId
}

func ReleaseAddress(svc *ec2.EC2, eipallocId string) error {

	input := &ec2.ReleaseAddressInput{
		AllocationId: aws.String(eipallocId),
	}

	result, err := svc.ReleaseAddress(input)
	//_, err := svc.ReleaseAddress(input)
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
		fmt.Printf("%s", result)
	}
	return err
}

func AssociateAddress(svc *ec2.EC2, instanceId, eipallocId string) string {

	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(eipallocId),
		InstanceId:   aws.String(instanceId),
	}

	result, err := svc.AssociateAddress(input)
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
		id = *result.AssociationId
	}
	if verbose {
		fmt.Println(result)
	}
	return id
}

func DisassociateAddress(svc *ec2.EC2, eipassocId string) error {

	input := &ec2.DisassociateAddressInput{
		AssociationId: aws.String(eipassocId),
	}

	result, err := svc.DisassociateAddress(input)
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

func PrintAddrInfo(v *ec2.Address) {
	var pip, alid, asid string
	pip = *v.PublicIp
	alid = *v.AllocationId
	if v.AssociationId != nil {
		asid = *v.AssociationId
	}
	tags := ""
	if v.Tags != nil {
		for _, tv := range v.Tags {
			tags += *tv.Key + " = " + *tv.Value + " | "
		}
	}
	fmt.Printf("%-20s %-30s %-30s %s\n", pip, alid, asid, tags)
}

// New stuff

func (i *Ec2Instance) AllocateAddress(svc *ec2.EC2) error {
	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain:  aws.String("vpc"),
		Address: aws.String(i.IP4),
		// Tagging after allocation
		/*
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String("elastic-ip"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String(ProjTagKey),
							Value: aws.String(ProjTagVal),
						},
					},
				},
			},
		*/
	})
	if allocRes.AllocationId != nil {
		//eipallocId = *allocRes.AllocationId
		i.IpAllocId = *allocRes.AllocationId
	}
	return err
}

func (e *Environment) ReleaseAddress(eipallocId string) error {

	input := &ec2.ReleaseAddressInput{
		AllocationId: aws.String(eipallocId),
	}

	result, err := e.SVC.ReleaseAddress(input)
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
		fmt.Printf("%s", result)
	}
	return err
}

func (i *Ec2Instance) AssociateAddress(svc *ec2.EC2) error {

	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(i.IpAllocId),
		InstanceId:   aws.String(i.InstanceID),
	}

	result, err := svc.AssociateAddress(input)
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
	if err == nil {
		//id = *result.AssociationId
		i.IpAssocId = *result.AssociationId
	}
	if verbose {
		fmt.Println(result)
	}
	return err
}

func (e *Environment) DisassociateAddress(eipassocId string) error {

	input := &ec2.DisassociateAddressInput{
		AssociationId: aws.String(eipassocId),
	}

	result, err := e.SVC.DisassociateAddress(input)
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

func (e *Environment) Ip4List() map[string]string {

	result, err := e.SVC.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	})

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

	ipList := make(map[string]string)
	if len(result.Addresses) > 0 {
		for _, v := range result.Addresses {

			if v.AssociationId != nil {
				ipList[*v.AllocationId] = *v.AssociationId
			} else {
				ipList[*v.AllocationId] = ""
			}
		}
	}

	return ipList
}

func init() {
	eipCMD.AddCommand(eipListCMD, eipDescribeCMD, eipAllocateCMD, eipReleaseCMD, eipAssociateCMD, eipDisassociateCMD)

	eipListCMD.Flags().StringVarP(&ip4, "ip4", "", "", "List IP address(es) based on IP rather than tags")
	eipAllocateCMD.Flags().StringVarP(&ip4, "ip4", "", "", "IP address")
	eipReleaseCMD.Flags().StringVarP(&eipallocId, "eipalloc-id", "", "", "EIP Allocation ID")
	eipAssociateCMD.Flags().StringVarP(&instanceId, "instance-id", "", "", "Instance ID")
	eipAssociateCMD.Flags().StringVarP(&eipallocId, "eipalloc-id", "", "", "EIP Allocation ID")
	eipDisassociateCMD.Flags().StringVarP(&eipassocId, "eipassoc-id", "", "", "EIP Association ID")
	eipDescribeCMD.Flags().StringVarP(&ip4, "ip4", "", "", "IP address")

	eipAllocateCMD.MarkFlagRequired("ip4")
	eipReleaseCMD.MarkFlagRequired("eipalloc-id")
	eipAssociateCMD.MarkFlagRequired("instance-id")
	eipAssociateCMD.MarkFlagRequired("eipalloc-id")
	eipDisassociateCMD.MarkFlagRequired("eipassoc-id")
	eipDescribeCMD.MarkFlagRequired("ip4")
}
