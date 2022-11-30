package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"
)

func yamlify() error {

	a := LabGroup{
		GroupName: "KTH01",
		Domain:    "kth01.examples.nu",
		Instances: []Instance{
			{
				Profile: "LAN",
				IP4:     "45.155.99.8",
				IP6:     "2a10:ba00:bee5::8",
			},
			{
				Profile: "Secondary",
				IP4:     "",
				IP6:     "",
			},
		},
	}
	b := LabGroup{
		GroupName: "KTH02",
		Domain:    "kth02.examples.nu",
		Instances: []Instance{
			{
				Profile: "LAN",
				IP4:     "45.155.99.12",
				IP6:     "2a10:ba00:bee5::12",
			},
			{
				Profile: "Secondary",
				IP4:     "",
				IP6:     "",
			},
		},
	}

	in := []LabGroup{a, b}

	out, _ := yaml.Marshal(&in)
	fmt.Printf("\n%v\n", string(out))

	return nil
}

func deyamlify(inFile string) {

	yamlFile, err := ioutil.ReadFile(inFile)
	var ret []LabGroup
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Printf("\n%+v\n", ret)
	for _, group := range ret {
		fmt.Printf("%s\n", group.GroupName)
		for _, i := range group.Instances {
			fmt.Printf("start an EC2 with profile %s\n", i.Profile)
		}
	}

}

func ReturnThis(a string) (error, string) {
	return nil, a
}

func AllocateByoipAddress(svc *ec2.EC2, ip4 string) (error, string) {
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

func AllocateAwsIpAddress(svc *ec2.EC2) (error, string) {
	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
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
