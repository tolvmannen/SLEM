package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2CMD = &cobra.Command{
	Use:   "ec2",
	Short: "Manually create / delete EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ec2 called\n")

	},
}

var ec2CreateCMD = &cobra.Command{
	Use:   "create",
	Short: "Manually create EC2 instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		// fix tags
		tags := []*ec2.Tag{
			{
				Key:   aws.String(ProjTagKey),
				Value: aws.String(ProjTagVal),
			},
		}

		err, instanceId = CreateEC2(svc, sgId, subnetId, tags)
		if err != nil {
			fmt.Printf("Failed to create EC2 Instance, %v", err)
		} else {
			fmt.Printf("EC2 Instance created: %s\n", instanceId)
		}

	},
}

var ec2DeleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Manually delete EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DeleteEC2(svc, instanceId)

	},
}

var ec2DescribeCMD = &cobra.Command{
	Use:   "describe",
	Short: "Describe EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = DescribeEC2(svc, instanceId)

	},
}

var ec2ListCMD = &cobra.Command{
	Use:   "list",
	Short: "List EC2 Instances",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		_ = ListEC2(svc)

	},
}

var ec2GetStateCMD = &cobra.Command{
	Use:   "state",
	Short: "Check state of an EC2 Instance",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		state := GetInstanceState(svc, instanceId)
		if state != "" {
			fmt.Printf("\nInstance %s is  %s\n\n", instanceId, state)
		}

	},
}

//Run
func CreateEC2(svc *ec2.EC2, sgId, subnetId string, tags []*ec2.Tag) (error, string) {

	input := &ec2.RunInstancesInput{
		//ImageId:      aws.String("ami-08edbb0e85d6a0a07"),
		ImageId:      aws.String("ami-0358232eacf458589"),
		InstanceType: aws.String("t3.micro"),
		KeyName:      aws.String("Jonas"),
		Ipv6Addresses: []*ec2.InstanceIpv6Address{
			{
				Ipv6Address: aws.String("2a10:ba00:bee5::8"),
			},
		},
		MaxCount: aws.Int64(1),
		MinCount: aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(sgId),
		},
		SubnetId: aws.String(subnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
				/*
					Tags: []*ec2.Tag{
						{
							Key:   aws.String(ProjTagKey),
							Value: aws.String(ProjTagVal),
						},
					},
				*/
			},
		},
	}

	result, err := svc.RunInstances(input)
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

	var instanceId string
	if err == nil {
		instanceId = *result.Instances[0].InstanceId
	}
	//instanceId := *result.Instances.InstanceId

	if verbose {
		fmt.Printf("%v\n", result)
	}
	//fmt.Printf("%T \n %#v", result, result)

	return err, instanceId

}

func CreateByoipEC2(svc *ec2.EC2, sgId, subnetId string, tags []*ec2.Tag) (error, string) {

	input := &ec2.RunInstancesInput{
		//ImageId:      aws.String("ami-08edbb0e85d6a0a07"),
		ImageId:      aws.String("ami-0358232eacf458589"),
		InstanceType: aws.String("t3.micro"),
		KeyName:      aws.String("Jonas"),
		Ipv6Addresses: []*ec2.InstanceIpv6Address{
			{
				Ipv6Address: aws.String("2a10:ba00:bee5::8"),
			},
		},
		MaxCount: aws.Int64(1),
		MinCount: aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(sgId),
		},
		SubnetId: aws.String(subnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
				/*
					Tags: []*ec2.Tag{
						{
							Key:   aws.String(ProjTagKey),
							Value: aws.String(ProjTagVal),
						},
					},
				*/
			},
		},
	}

	result, err := svc.RunInstances(input)
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

	var instanceId string
	if err == nil {
		instanceId = *result.Instances[0].InstanceId
	}
	//instanceId := *result.Instances.InstanceId

	if verbose {
		fmt.Printf("%v\n", result)
	}
	//fmt.Printf("%T \n %#v", result, result)

	return err, instanceId

}

func CreateAwsIpEC2(svc *ec2.EC2, sgId, subnetId string, tags []*ec2.Tag) (error, string) {

	input := &ec2.RunInstancesInput{
		ImageId:          aws.String("ami-0358232eacf458589"),
		InstanceType:     aws.String("t3.micro"),
		KeyName:          aws.String("Jonas"),
		Ipv6AddressCount: aws.Int64(1),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(sgId),
		},
		SubnetId: aws.String(subnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
			},
		},
	}

	result, err := svc.RunInstances(input)
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

	var instanceId string
	if err == nil {
		instanceId = *result.Instances[0].InstanceId
	}
	//instanceId := *result.Instances.InstanceId

	if verbose {
		fmt.Printf("%v\n", result)
	}
	//fmt.Printf("%T \n %#v", result, result)

	return err, instanceId

}

// Terminate
func DeleteEC2(svc *ec2.EC2, instanceId string) error {

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.TerminateInstances(input)
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

func DescribeEC2(svc *ec2.EC2, instanceId string) error {

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.DescribeInstances(input)
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
	var id, state string
	if len(result.Reservations) > 0 {
		id = *result.Reservations[0].Instances[0].InstanceId
		state = *result.Reservations[0].Instances[0].State.Name
	}
	fmt.Printf("\nInstande ID: %s\nState: %s\n", id, state)

	return nil
}

func ListEC2(svc *ec2.EC2) error {

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + ProjTagKey),
				Values: []*string{
					aws.String(ProjTagVal),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
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

	//var iId, state, pubIp string
	fmt.Printf("%-20s %-20s %-30s %-20s %-30s %s\n", "Group name", "Public IP", "Public IPv6", "State", "InstanceId", "Tags")
	for _, ec2inst := range result.Reservations {
		iId := *ec2inst.Instances[0].InstanceId
		state := *ec2inst.Instances[0].State.Name
		pubIp := ""
		if ec2inst.Instances[0].PublicIpAddress != nil {
			pubIp = *ec2inst.Instances[0].PublicIpAddress
		}
		pubIp6 := ""
		if ec2inst.Instances[0].Ipv6Address != nil {
			pubIp6 = *ec2inst.Instances[0].Ipv6Address

		}
		tags := ""
		gName := ""
		if ec2inst.Instances[0].Tags != nil {
			for _, tv := range ec2inst.Instances[0].Tags {
				tags += *tv.Key + " = " + *tv.Value + " | "
				if *tv.Key == "Name" {
					gName = *tv.Value
				}
			}
		}
		fmt.Printf("%-20s %-20s %-30s %-20s %-30s %s\n", gName, pubIp, pubIp6, state, iId, tags)
		//fmt.Printf("%T\n%+v\n-----------\n", ec2inst, ec2inst)
		//fmt.Printf("%+v\n-----------\n", ec2inst.Instances[0].InstanceId)

	}
	/*
		fmt.Println(result)
		var id, state string
		if len(result.Reservations) > 0 {
			id = *result.Reservations[0].Instances[0].InstanceId
			state = *result.Reservations[0].Instances[0].State.Name
		}
		fmt.Printf("\nInstande ID: %s\nState: %s\n", id, state)
	*/

	return nil
}

func GetInstanceId(svc *ec2.EC2, instanceId string) string {

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.DescribeInstances(input)
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
	if len(result.Reservations) > 0 {
		id = *result.Reservations[0].Instances[0].InstanceId
	}
	//fmt.Printf("%T \n %#v", result, result)
	return id
}

func GetInstanceState(svc *ec2.EC2, instanceId string) string {

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.DescribeInstances(input)
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

	var state string
	if len(result.Reservations) > 0 {
		state = *result.Reservations[0].Instances[0].State.Name
	}
	//fmt.Printf("\nInstande ID: %s\nState: %s\n", instanceId, state)

	return state
}

func init() {
	ec2CMD.AddCommand(ec2CreateCMD, ec2DeleteCMD, ec2DescribeCMD, ec2GetStateCMD, ec2ListCMD)
	ec2CreateCMD.Flags().StringVarP(&sgId, "sg-id", "", "", "Security Group ID")
	ec2CreateCMD.Flags().StringVarP(&subnetId, "subnet-id", "", "", "Subnet ID")
	ec2DeleteCMD.Flags().StringVarP(&instanceId, "instance-id", "", "", "Instance ID")
	ec2DescribeCMD.Flags().StringVarP(&instanceId, "instance-id", "", "", "Instance ID")
	ec2GetStateCMD.Flags().StringVarP(&instanceId, "instance-id", "", "", "Instance ID")

	ec2CreateCMD.MarkFlagRequired("sg-id")
	ec2CreateCMD.MarkFlagRequired("subnet-id")
	ec2DeleteCMD.MarkFlagRequired("instance-id")
	ec2DescribeCMD.MarkFlagRequired("instance-id")
	ec2GetStateCMD.MarkFlagRequired("instance-id")
}
