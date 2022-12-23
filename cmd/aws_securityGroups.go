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

		e.CreateSG()
		e.AddIngressRules()

	},
}

var sgDeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete Security Group",
	Run: func(cmd *cobra.Command, args []string) {

		if e.SgId == "" {
			e.SgId = sgId
		}
		e.DeleteSG()
	},
}

var sgDescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe Security Group",
	Run: func(cmd *cobra.Command, args []string) {

		e.DescribeSG(sgId)

	},
}

var sgListCMD = &cobra.Command{
	Use:   "list",
	Short: "List Security Groups by tags",
	Run: func(cmd *cobra.Command, args []string) {

		e.ListSG()

	},
}

// Bundling functions

//func (e *Environment) GetSgId(svc *ec2.EC2) error {
func (e *Environment) GetSgId() error {

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:EnvTag"),
				Values: []*string{
					aws.String(e.EnvTag),
				},
			},
		},
	}

	result, err := e.SVC.DescribeSecurityGroups(input)
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

	for _, sgs := range result.SecurityGroups {
		e.SgId = *sgs.GroupId
	}

	if len(result.SecurityGroups) < 1 {
		err = noIdErr
	} else if len(result.SecurityGroups) > 1 {
		err = xsIdErr
	}

	return err

}

func (e *Environment) CreateSG() error {

	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String(e.EnvTag + "-SG"),
		GroupName:   aws.String(e.EnvTag + "-SG"),
		VpcId:       aws.String(e.VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("EnvTag"),
						Value: aws.String(e.EnvTag),
					},
				},
			},
		},
	}

	result, err := e.SVC.CreateSecurityGroup(input)
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

	if result.GroupId != nil {
		e.SgId = *result.GroupId
	}

	if verbose {
		fmt.Println(result)
	}
	return err

}

func (e *Environment) AddIngressRules() error {

	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       aws.String(e.SgId),
		IpPermissions: EC2LanRules,
	}

	result, err := e.SVC.AuthorizeSecurityGroupIngress(input)
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

func (e *Environment) DeleteSG() error {

	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(e.SgId),
		//GroupId: aws.String(sgId),
	}

	result, err := e.SVC.DeleteSecurityGroup(input)
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

func (e *Environment) DescribeSG(sgId string) error {

	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(sgId),
		},
	}

	result, err := e.SVC.DescribeSecurityGroups(input)
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

func (e *Environment) ListSG() error {

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	result, err := e.SVC.DescribeSecurityGroups(input)
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

func init() {
	sgCMD.AddCommand(sgCreateCMD, sgDeleteCMD, sgDescribeCMD, sgListCMD)
	sgCreateCMD.Flags().StringVarP(&vpcId, "vpc-id", "", "", "VPC ID")
	sgDeleteCMD.Flags().StringVarP(&sgId, "sg-id", "", "", "Security Group ID")
	sgDescribeCMD.Flags().StringVarP(&sgId, "sg-id", "", "", "Security Group ID")
	sgCreateCMD.MarkFlagRequired("vpc-id")
	sgDeleteCMD.MarkFlagRequired("sg-id")
	sgDescribeCMD.MarkFlagRequired("sg-id")
}
