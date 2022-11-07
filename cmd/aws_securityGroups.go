package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var sgCMD = &cobra.Command{
	Use:   "sg",
	Short: "Manually create / delete Security Group",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sg called\n")

	},
}

var sgCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create Security Group",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		err, sgId := CreateSG(svc, vpcId)
		if err == nil {
			fmt.Printf("SG created: %s\n", sgId)
			err = AddIngressRulesLAN(svc, sgId)
			if err == nil {
				fmt.Printf("Ingress Rules added to Security Group %s\n", err)
			}
		}

	},
}

var sgDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Security Group",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteSG(svc, sgId)

	},
}

var sgDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Security Group",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeSG(svc, sgId)

	},
}

var sgListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Security Groups by tags",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = ListSG(svc)

	},
}

func CreateSG(svc *ec2.EC2, vpcId string) (error, string) {

	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String("DNS-course LAN SG"),
		GroupName:   aws.String("DNS-course-LAN-SG"),
		VpcId:       aws.String(vpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String(ProjTagKey),
						Value: aws.String(ProjTagVal),
					},
				},
			},
		},
	}

	result, err := svc.CreateSecurityGroup(input)
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

	var gid string
	if result.GroupId != nil {
		gid = *result.GroupId
	}

	if verbose {
		fmt.Println(result)
	}
	return err, gid

}

func DeleteSG(svc *ec2.EC2, sgId string) error {

	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgId),
	}

	result, err := svc.DeleteSecurityGroup(input)
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

func DescribeSG(svc *ec2.EC2, sgId string) error {

	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(sgId),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
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

func ListSG(svc *ec2.EC2) error {

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
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

	//fmt.Println(result)
	fmt.Printf("%-30s %-30s %s\n", "SG-ID", "vpc-id", "Tags")

	for _, sgs := range result.SecurityGroups {
		sgId := *sgs.GroupId
		vpcId := *sgs.VpcId
		tags := ""
		if sgs.Tags != nil {
			for _, tv := range sgs.Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
			}
		}
		fmt.Printf("%-30s %-30s %s\n", sgId, vpcId, tags)
	}

	return nil

}

func AddIngressRulesLAN(svc *ec2.EC2, sgId string) error {

	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(sgId),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("SSH access to Docker"),
					},
				},
				ToPort: aws.Int64(22),
			},
			{
				FromPort:   aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				Ipv6Ranges: []*ec2.Ipv6Range{
					{
						CidrIpv6:    aws.String("::0/0"),
						Description: aws.String("SSH to Docker IPv6"),
					},
				},
				ToPort: aws.Int64(22),
			},
			{
				FromPort:   aws.Int64(25),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("SMTP access"),
					},
				},
				ToPort: aws.Int64(25),
			},
			{
				FromPort:   aws.Int64(53),
				IpProtocol: aws.String("udp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("DNS"),
					},
				},
				ToPort: aws.Int64(53),
			},
			{
				FromPort:   aws.Int64(53),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("DNS"),
					},
				},
				ToPort: aws.Int64(53),
			},
			{
				FromPort:   aws.Int64(53),
				IpProtocol: aws.String("udp"),
				Ipv6Ranges: []*ec2.Ipv6Range{
					{
						CidrIpv6:    aws.String("::0/0"),
						Description: aws.String("DNS IPv6"),
					},
				},
				ToPort: aws.Int64(53),
			},
			{
				FromPort:   aws.Int64(53),
				IpProtocol: aws.String("tcp"),
				Ipv6Ranges: []*ec2.Ipv6Range{
					{
						CidrIpv6:    aws.String("::0/0"),
						Description: aws.String("DNS IPv6"),
					},
				},
				ToPort: aws.Int64(53),
			},
			{
				FromPort:   aws.Int64(80),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("HTTP"),
					},
				},
				ToPort: aws.Int64(80),
			},
			{
				FromPort:   aws.Int64(80),
				IpProtocol: aws.String("tcp"),
				Ipv6Ranges: []*ec2.Ipv6Range{
					{
						CidrIpv6:    aws.String("::0/0"),
						Description: aws.String("HTTP IPv6"),
					},
				},
				ToPort: aws.Int64(80),
			},
			{
				FromPort:   aws.Int64(2022),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("SSH access to EC2"),
					},
				},
				ToPort: aws.Int64(2022),
			},
		},
	}

	result, err := svc.AuthorizeSecurityGroupIngress(input)
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
		fmt.Printf("%v", result)
	}
	return err

}

func init() {
	sgCMD.AddCommand(sgCreateCMD, sgDeleteCMD, sgDescribeCMD, sgListCMD)
	sgCreateCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	sgDeleteCMD.Flags().StringVarP(&sgId, "sg-id", "", "", "Security Group ID")
	sgDescribeCMD.Flags().StringVarP(&sgId, "sg-id", "", "", "Security Group ID")
	sgCreateCMD.MarkFlagRequired("vpc-id")
	sgDeleteCMD.MarkFlagRequired("sg-id")
	sgDescribeCMD.MarkFlagRequired("sg-id")
}
