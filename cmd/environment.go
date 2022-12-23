package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"
)

type Environment struct {
	Project string `yaml:"Project"`
	EnvTag  string `yaml:"EnvTag"`
	Type    string `yaml:"Type"`

	AWSProfile string `yaml:"AWSProfile"`
	AWSRegion  string `yaml:"AWSRegion"`
	AWSKeypair string `yaml:"AWSKeypair"`

	Byoip4cidr   string `yaml:"Byoip4cidr"`
	Byoip6cidr   string `yaml:"Byoip6cidr"`
	Byoip6pool   string `yaml:"Byoip6pool"`
	ImageId      string `yaml:"ImageId"`
	InstanceType string `yaml:"InstanceType"`
	//Domain       string `yaml:"Domain"`
	StartIP     int `yaml:"StartIP"`
	GroupIpSize int `yaml:"GroupIpSize"`
	NrOfGroups  int `yaml:"NrOfGroups"`

	// updated during install
	VpcId      string `yaml:"VpcId"`
	IgwId      string `yaml:"IgwId"`
	SubnetId   string `yaml:"SubnetId"`
	RtbId      string `yaml:"RtbId"`
	RtbassocId string `yaml:"RtbassocId"`
	SgId       string `yaml:"SgId"`

	SVC *ec2.EC2

	Ec2Instances []Ec2Instance `yaml:"Ec2Instances"`
}

type Ec2Instance struct {
	InstanceID string `yaml:"InstanceID"`
	IpAllocId  string `yaml:"IpAllocId"`
	IpAssocId  string `yaml:"IpAssocId"`
	IP4        string `yaml:"IP4"`
	IP6        string `yaml:"IP6"`
	State      string `yaml:"State"`
}

func (e *Environment) LoadConf(inFile string) error {

	yamlFile, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &e)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return err
}

func (e *Environment) CreateSession() error {

	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: e.AWSProfile,

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(e.AWSRegion),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	e.SVC = ec2.New(sess)

	return err
}

func (e *Environment) Load() error {
	err = e.GetVpcId()
	err = e.GetIgwId()
	err = e.GetSubnetId()
	err = e.GetRtbId()
	err = e.GetRtbassocId()
	err = e.GetSgId()
	e.Ec2List()

	return err
}

func (e *Environment) MakeIpPlan() error {

	ip := e.StartIP
	ammount := e.NrOfGroups
	inc := e.GroupIpSize

	// # Only provision Instances in exess of what is already up and running
	current := len(e.Ec2Instances)
	fmt.Printf("Instances:\n\t%-20s(%v)\n\t%-20s(%v)\n", "Target:", ammount, "Running:", current)
	if ammount <= current {
		exitErrorf("Targen nr reached (or exceeded). Skipping...\n")
	} else {
		ammount = ammount - current
		fmt.Printf("Creating  an additional %v Instances\n", ammount)
	}

	for i := 0; i < ammount; i++ {
		if ip > 255 {
			exitErrorf("\nIP out of bounds! (%v) Panic!\n", ip)
		}
		p := Ec2Instance{
			IP4: "45.155.99." + strconv.Itoa(ip),
			IP6: "2a10:ba00:bee5::" + strconv.Itoa(ip),
		}
		err, id := e.Ip4Status(p.IP4)
		// ## If IP allocated (but not associated), add allocation id to instance info
		// ## or allocatiopn will (soft) fail further on
		if !errors.Is(err, ipAssErr) && !e.Ip6InUse(p.IP6) {
			if errors.Is(err, ipAllErr) {
				p.IpAllocId = id
			}
			e.Ec2Instances = append(e.Ec2Instances, p)
		} else {
			ammount++
		}
		ip += inc
	}

	return err
}

func (e *Environment) Ip4Status(ip4 string) (error, string) {
	iplist := e.GetAllocatedV4(ip4)
	var id string
	//fmt.Printf("\n%v\n", iplist)
	if iplist != nil {
		//if iplist[0].AllocationId != nil {
		if iplist[0].AssociationId != nil {
			id = *iplist[0].AssociationId
			err = ipAssErr
		} else {
			id = *iplist[0].AllocationId
			err = ipAllErr
		}
	}

	return err, id
}

func (e *Environment) GetAllocatedV4(ip4 string) []*ec2.Address {

	result, err := e.SVC.DescribeAddresses(&ec2.DescribeAddressesInput{
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

	ret := result.Addresses
	return ret
}

// Move to eip?
func (e *Environment) Ip6InUse(ip6 string) bool {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Project"),
				Values: []*string{
					aws.String(e.Project),
				},
			},
		},
	}

	result, err := e.SVC.DescribeInstances(input)
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

	for _, ec2inst := range result.Reservations {
		state := *ec2inst.Instances[0].State.Name

		if state != "terminated" {
			/*
				pubIp := ""
				if ec2inst.Instances[0].PublicIpAddress != nil {
					pubIp = *ec2inst.Instances[0].PublicIpAddress
				}
				pubIp6 := ""
				if ec2inst.Instances[0].Ipv6Address != nil {
					pubIp6 = *ec2inst.Instances[0].Ipv6Address

				}
			*/
			if *ec2inst.Instances[0].Ipv6Address == ip6 {
				return true
			}
		}

	}

	return false
}

func (e *Environment) DeployEC2() error {

	for k, i := range e.Ec2Instances {
		err = e.CreateByoipEC2(&i)
		if err != nil {
			fmt.Printf("IP in use?\n")
		} else {
			if i.IpAllocId == "" {
				i.AllocateAddress(e.SVC)
				AddTag(e.SVC, i.IpAllocId, "EnvTag", e.EnvTag)
				AddTag(e.SVC, i.IpAllocId, "Project", e.Project)
			} else {
				fmt.Printf("%s already allocated\n", i.IP4)
			}
			state := GetInstanceState(e.SVC, i.InstanceID)
			tries := 0
			for {
				if state != "running" {
					if tries > 9 {
						exitErrorf("Timeout waiting for EC2 to start")
					}
					time.Sleep(1 * time.Second)
					state = GetInstanceState(e.SVC, i.InstanceID)
					tries++
				} else {
					break
				}

				i.AssociateAddress(e.SVC)

				e.Ec2Instances[k] = i
			}
		}
	}

	return err
}

func (e *Environment) DestroyEC2() error {

	//#  Work order:
	//#	1. Disassociate addresses
	//#	2. Release addresses
	//#	3. Delete EC2 Instances

	ip4list := e.Ip4List()
	for alId, asId := range ip4list {
		if asId != "" {
			e.DisassociateAddress(asId)
		}
		e.ReleaseAddress(alId)
	}

	ec2list := e.Ec2List()
	for _, ec2id := range ec2list {
		e.DeleteEC2(ec2id)
	}
	// # termionating an instance may take a little bit
	// # keep checking every 5 seconds for a minute
	// # then terminate.
	tries := 0
	for {
		if len(ec2list) > 0 {
			if tries > 24 {
				exitErrorf("Timeout waiting for all instances to terminate\n")
			}
			fmt.Printf("Waiting for instances to terminate: (%v left)\n", len(ec2list))
			time.Sleep(5 * time.Second)
			ec2list = e.Ec2List()
			tries++
		} else {
			break
		}
	}

	return err
}

// 1. DEPLOY
func (e *Environment) DeployEnv() error {

	// === VPC ===
	err = e.GetVpcId()
	if errors.Is(err, noIdErr) {
		switch e.Type {
		case "BYOIP":
			fmt.Printf("\nCreating VPC  (%s)\n", e.EnvTag)
			//e.CreateVPC(c.SVC)
			e.CreateByoipVpc()
			AddTag(e.SVC, e.VpcId, "Project", e.Project)
		default:
			exitErrorf("\nUNKNOWN TYPE (%v) Aborting!\n", e.Type)
		}
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		e.DescribeVpc(e.VpcId)
		fmt.Printf("\n\nFound VPC (id: %s)!\n", e.VpcId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}

	// === IGW ===
	err = e.GetIgwId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nCreating Internet Gateway  (%s)\n", e.EnvTag)
		e.CreateIGW()
		AddTag(e.SVC, e.IgwId, "Project", e.Project)
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		e.DescribeIGW(e.IgwId)
		fmt.Printf("\n\nFound IGW (id: %s)!\n", e.IgwId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}
	// === IGW - attach ===
	// Attach IGW. If already attached, nothing happens
	fmt.Printf("\nAttachin Internet Gateway (%s) to VPC (%s)\n", e.IgwId, e.VpcId)
	err = e.AttachIGW()

	// === Subnet ===
	err = e.GetSubnetId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nCreating Subnet (%s)\n", e.EnvTag)
		e.CreateSubnet()
		AddTag(e.SVC, e.SubnetId, "Project", e.Project)
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		e.DescribeSubnet(e.SubnetId)
		fmt.Printf("\n\nFound Subnet (id: %s)!\n", e.SubnetId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}

	// === Routing Table ===
	err = e.GetRtbId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nCreating Routing Table (%s)\n", e.EnvTag)
		e.CreateRTB()
		AddTag(e.SVC, e.RtbId, "Project", e.Project)
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		//_ = DescribeRT(c.SVC, e.RtbId)
		e.DescribeRTB(e.RtbId)
		fmt.Printf("\n\nFound Route Table (id: %s)!\n", e.RtbId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}

	// === Routing Table - association ===
	// Attach RTB. If already attached, nothing happens
	err = e.GetRtbassocId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nAssociating Routing Table (%s) to Subnet (%s)\n", e.RtbId, e.SubnetId)
		e.AssociateRTB()
		//AddTag(c.SVC, e.RtbassocId, "Project", c.Project)
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		e.DescribeRTB(e.RtbId)
		fmt.Printf("\n\nAssociation to Route Table (id: %s)!\n", e.RtbId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}

	// === Route v4 ===
	// If already created, nothing happens
	fmt.Printf("Creating IPv4 Route attached to IGW (%s) and RTB (%s)\n", e.IgwId, e.RtbId)
	err = e.CreateRoute()

	// === Route v6 ===
	// If already created, nothing happens
	fmt.Printf("Creating IPv6 Route attached to IGW (%s) and RTB (%s)\n", e.IgwId, e.RtbId)
	err = e.CreateRoute6()

	// === Security Group ===
	//fmt.Printf("\nSG RULES\n%+v\n", EC2LanRules)
	err = e.GetSgId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nCreating Security Group (%s)\n", e.EnvTag)
		err = e.CreateSG()
		AddTag(e.SVC, e.SgId, "Project", e.Project)
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		var cont string
		e.DescribeSG(e.SgId)
		fmt.Printf("\n\nFound Security Group (id: %s)!\n", e.SgId)
		fmt.Printf("\nEnter \"y\" to continue, or anything else to abort.\n")
		fmt.Scanf("%s", &cont)
		if cont != "y" {
			exitErrorf("\nAborting!\n")
		}
	}

	// === Security Group - Ingress Rules ===
	fmt.Printf("\nAdding ingress rules to Security Group (%s)\n", e.SgId)
	err = e.AddIngressRules()

	//c.Environment[k] = e

	return err
}

// 2. DESTROY E
func (e *Environment) DestroyEnv() error {
	// === Security Group ===
	if e.SgId != "" {
		fmt.Printf("\nDeleting Security Group (id: %s)!\n", e.SgId)
		e.DeleteSG()
	} else {
		fmt.Printf("\nNo Security Group found. Skipping...\n")
	}

	// === Routing Table - association ===
	if e.RtbassocId != "" {
		fmt.Printf("\nDisassociating Route Table (RTB Association ID: %s)\n", e.RtbassocId)
		e.DisassociateRTB()
	} else {
		fmt.Printf("\nNo RTB association found. Skipping...\n")
	}

	// === Routing Table ===
	if e.RtbId != "" {
		fmt.Printf("\nDeleting Route Table (id: %s)\n", e.RtbId)
		e.DeleteRTB()
	} else {
		fmt.Printf("\nNo Routing Table found. Skipping...\n")
	}

	// === Subnet ===
	if e.SubnetId != "" {
		fmt.Printf("\nDeleting Subnet (id: %s)\n", e.SubnetId)
		e.DeleteSubnet()
	} else {
		fmt.Printf("\nNo Subnet found. Skipping...\n")
	}

	// === IGW - detach ===
	// # Nothing happens if not attached..
	if e.IgwId != "" && e.VpcId != "" {
		fmt.Printf("\nDetatching Internet Gateway (id: %s / %)\n", e.IgwId, e.VpcId)
		e.DetatchIGW()
	}

	// === IGW ===
	if e.IgwId != "" {
		fmt.Printf("\nDeleting Internet Gateway (id: %s)\n", e.IgwId)
		e.DeleteIGW()
	} else {
		fmt.Printf("\nNo Internet Gateway  found. Skipping...\n")
	}

	// === VPC ===
	if e.VpcId != "" {
		fmt.Printf("\nDeleting Internet Gateway (id: %s)\n", e.VpcId)
		e.DeleteVpc()
	} else {
		fmt.Printf("\nNo VPC found. Skipping...\n")
	}

	return err
}

/*
func (e *Environment) DestroyEnv() error {

	// === Security Group ===
	err = e.GetSgId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nNo Security Group found. Skipping...\n")
	} else if errors.Is(err, exIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDeleting Security Group (id: %s)!\n", e.SgId)
		e.DeleteSG(e.SgId)
	}

	// === Routing Table - association ===
	err = e.GetRtbassocId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nNo RTB association found. Skipping...\n")
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDisassociating Route Table (RTB Association ID: %s)\n", e.RtbassocId)
		e.DisassociateRTB()
	}

	// === Routing Table ===
	err = e.GetRtbId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nNo RTB association found. Skipping...\n")
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDeleting Route Table (id: %s)!\n", e.RtbId)
		e.DeleteRTB()
	}

	// === Subnet ===
	err = e.GetSubnetId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nNo Subnet found. Skipping...\n")
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDeleting Subnet (%s)\n", e.SubnetId)
		e.DeleteSubnet()
	}

	// === IGW - detach ===
	fmt.Printf("\nDetachin Internet Gateway (%s) from  VPC (%s)\n", e.IgwId, e.VpcId)
	// === IGW ===
	err = e.GetIgwId()
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nNo Creating Internet Gateway found. Skipping...\n")
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDeletingn Internet Gateway (%s)\n", e.IgwId)
		var cont string
		e.DeleteIGW()
	}

	// === VPC ===
	err = e.GetVpcId(e.SVC)
	if errors.Is(err, noIdErr) {
		fmt.Printf("\nVPC found. Skipping...\n")
	} else if errors.Is(err, xsIdErr) {
		exitErrorf("%v", err)
	} else {
		fmt.Printf("\nDeleting VPC  (%s)\n", e.VpcId)
		e.DeleteVpc()
	}

	return err
}
*/
