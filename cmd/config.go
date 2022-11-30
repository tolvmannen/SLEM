package cmd

import (
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"
)

type SlemConfig struct {
	// DEFAULTS - Added from main config
	AWSProfile   string              `yaml:"AWSProfile"`
	AWSRegion    string              `yaml:"AWSRegion"`
	AWSKeypair   string              `yaml:"AWSKeypair"`
	Byoip4cidr   string              `yaml:"Byoip4cidr"`
	Byoip6cidr   string              `yaml:"Byoip6cidr"`
	ImageId      string              `yaml:"ImageId"`
	InstanceType string              `yaml:"InstanceType"`
	ProjectName  string              `yaml:"ProjectName"`
	Tags         []map[string]string `yaml:"Tags"`
	// Added from separate config
	LabGroups []LabGroup
	// Updated during deployment
	VpcId      string
	IgwId      string
	SubnetId   string
	RtbId      string
	RtbassocId string
	SgId       string
}

type LabGroup struct {
	GroupName string     `yaml:"GroupName"`
	Domain    string     `yaml:"Domain"`
	Instances []Instance `yaml:"Instance"`
}

type Instance struct {
	Profile string `yaml:"Profile"`
	ImageId string `yaml:"ImageId"`
	IP4     string `yaml:"IP4"`
	IP6     string `yaml:"IP6"`

	EipassocId string
	EipallocId string
	InstanceId string
}

// not used yet
type EnvConf struct {
	IgwId      string
	VpcId      string
	SubnetId   string
	SgId       string
	RtbId      string
	RtbassocId string
}

func LoadSlemConf(inFile string) (SlemConfig, error) {

	yamlFile, err := ioutil.ReadFile(inFile)
	var ret SlemConfig
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ret, err

}

func ApendTags(target []*ec2.Tag, tags []map[string]string) []*ec2.Tag {
	for _, ctag := range tags {
		for k, v := range ctag {
			//fmt.Printf("Values in: %s : %s\n", k, v)
			var tag ec2.Tag
			tag = ec2.Tag{
				Key:   aws.String(k),
				Value: aws.String(v),
			}
			target = append(target, &tag)
		}
	}
	return target

}
