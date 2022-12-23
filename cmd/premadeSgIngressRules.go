package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var EC2LanRules = []*ec2.IpPermission{
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
}
