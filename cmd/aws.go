package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	ProjTagKey = "Project"
	ProjTagVal = "DNS-course"
)

var (
	ipstart int
	iprange int

	awsconf string
	labconf string

	igwId      string
	vpcId      string
	subnetId   string
	sgId       string
	rtbId      string
	rtbassocId string
	eipassocId string
	eipallocId string
	instanceId string
	eniId      string
	ip4        string
	ip6        string

	tagName    string
	tagValue   string
	resourceId string

	sess2 *session.Session
	svc2  *ec2.EC2
	err   error
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Manage the AWS infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws called")
	},
}

var testCMD = &cobra.Command{
	Use:   "test",
	Short: "Test stuff",
	Run: func(cmd *cobra.Command, args []string) {

		SlemConf, _ := LoadSlemConf(awsconf)
		fmt.Printf("SlemConf = %+v\n", SlemConf)

		var tags []*ec2.Tag
		tags = ApendTags(tags, SlemConf.Tags)

		/*
			extratags := []map[string]string{
				{
					"Key":   "Test",
					"Value": "Successful?",
				},
			}
			tags = ApendTags(tags, extratags)
		*/

		//fmt.Printf("%T : %#v\n", tags, tags)
		//inFile := "./config/test.yaml"

		//SlemConf.LabGroups = LoadLabGroupConf(inFile)
		//fmt.Printf("SlemConf = %+v\n", SlemConf)
		/*
			fmt.Printf("Creating VPC\n")
			_, SlemConf.VpcId = ReturnThis("vpc-0f152cbe7449a4e37")
			// _, vpcId = CreateVpc(svc, Tags)
			fmt.Printf("Creating Internet Gateway\n")
			_, SlemConf.IgwId = ReturnThis("igw-0f152cbe7449a4e99")
			// _, igwId = CreateInternetGateway(svc, Tags)
			fmt.Printf("Attaching Internet Gateway to VPC\n")
			// _ = AttachInternetGateway(svc, igwId, vpcId)
			fmt.Printf("Creating Subnet\n")
			_, SlemConf.SubnetId = ReturnThis("subnet-0f152cbe77346s99")
			// _, subnetId := CreateSubnet(svc, vpcId, byoip4, byoip6. Tags)
			fmt.Printf("Creating Routing Table\n")
			_, SlemConf.RtbId = ReturnThis("rtb-0f152cbe77346s77")
			// _, rtbId = CreateRT(svc, vpcId, Tags)
			fmt.Printf("Associating Routing Table with Subnet\n")
			_, SlemConf.RtbassocId = ReturnThis("rtbassoc-0f152cbe7449a4e99")
			// _, rtbassocId = AssociateRT(svc, rtbId, subnetId)
			fmt.Printf("Creating IPv4 Route attached to IGW and Routnig Table\n")
			// _ =  = CreateRoute(svc, igwId, rtbId)
			fmt.Printf("Creating IPv6 Route attached to IGW and Routnig Table\n")
			// _ = CreateRoute6(svc, igwId, rtbId)
			fmt.Printf("Creating Security Group\n")
			_, SlemConf.SgId = ReturnThis("sg-0f152cbe7449a4e37")
			// err, sgId := CreateSG(svc, vpcId, Tags)
			fmt.Printf("Adding Ingress Rules to Security Group\n")
			// _ = AddIngressRulesLAN(svc, sgId)

			fmt.Printf("\n=======================\n")
			fmt.Printf("Ney SlemConf = %+v\n", SlemConf)
		*/
		//LabGroups := LoadLabGroupConf(inFile)
		/*
			for _, g := range LabGroups {
				fmt.Printf("%+v\n", g)
				for _, i := range g.Instances {
					if i.ImageId == "" {
						//fmt.Printf("using default AMI: %s\n", SlemConf.ImageId)
						//fmt.Printf("using SPESHUL AMI: %s\n", i.ImageId)
						i.ImageId = SlemConf.ImageId
					}
					//if i.ImageId == "" {
					//}
					fmt.Printf("using AMI: %s\n", i.ImageId)

				}
			}

			// ----
			sess, err := CreateAwsSession()

			if err != nil {
				fmt.Printf("Session create error, %v", err)
			}

			svc := ec2.New(sess)
			if ProjectVpcExists(svc) {
				fmt.Printf("\nVPC already exists!\n")
			}
		*/
		// ---

		//deyamlify(inFile)
		//yamlify()
	},
}

func DeployEnvironment(svc *ec2.EC2, Tags []*ec2.Tag) error {
	//fmt.Printf("Maybe make wrapper later...\n")

	DeployNetconf(svc, Tags)

	return nil
}

func DeployNetconf(svc *ec2.EC2, Tags []*ec2.Tag) error {

	fmt.Printf("Creating VPC\n")
	err, vpcId = CreateVpc(svc)
	// err, vpcId = CreateVpc(svc, Tags)
	if err != nil {
		exitErrorf("VPC exists")
	}
	fmt.Printf("Creating Internet Gateway\n")
	// _, igwId = CreateInternetGateway(svc, Tags)
	fmt.Printf("Attaching Internet Gateway to VPC\n")
	// _ = AttachInternetGateway(svc, igwId, vpcId)
	fmt.Printf("Creating Subnet\n")
	// _, subnetId := CreateSubnet(svc, vpcId, byoip4, byoip6. Tags)
	fmt.Printf("Creating Routing Table\n")
	// _, rtbId = CreateRT(svc, vpcId, Tags)
	fmt.Printf("Associating Routing Table with Subnet\n")
	// _, rtbassocId = AssociateRT(svc, rtbId, subnetId)
	fmt.Printf("Creating IPv4 Route attached to IGW and Routnig Table\n")
	// _ =  = CreateRoute(svc, igwId, rtbId)
	fmt.Printf("Creating IPv6 Route attached to IGW and Routnig Table\n")
	// _ = CreateRoute6(svc, igwId, rtbId)
	fmt.Printf("Creating Security Group\n")
	// err, sgId := CreateSG(svc, vpcId, Tags)
	fmt.Printf("Adding Ingress Rules to Security Group\n")
	// _ = AddIngressRulesLAN(svc, sgId)

	return nil

}

/*
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

*/
func LoadLabGroupConf(inFile string) []LabGroup {

	yamlFile, err := ioutil.ReadFile(inFile)
	var ret []LabGroup
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return ret

}

var addtagCMD = &cobra.Command{
	Use:   "tag",
	Short: "Tag a reource",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		svc := ec2.New(sess)

		AddTag(svc, resourceId, tagName, tagValue)
	},
}

var untagCMD = &cobra.Command{
	Use:   "untag",
	Short: "Remove a tag from a reource",
	Run: func(cmd *cobra.Command, args []string) {

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		svc := ec2.New(sess)

		UnTag(svc, resourceId, tagName, tagValue)
	},
}

var ipcalcCMD = &cobra.Command{
	Use:   "ipcalc",
	Short: "Test IP addressing",
	Run: func(cmd *cobra.Command, args []string) {
		var ip4, ip6 string
		for i := ipstart; i < 255; i += iprange {
			ip4 = "45.155.99." + strconv.Itoa(i)
			lh := fmt.Sprintf("%04x", i)
			//ip6 = "2a10:ba00:bee5:" + lh + ":0000:0000:0000:0022"
			ip6 = "2a10:ba00:bee5::" + lh
			fmt.Printf("\nIPv4: %s\nIPv6: %s\n", ip4, ip6)
		}
	},
}

func CreateAwsSession() (*session.Session, error) {

	// Initialize a session in <region> that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: "default",

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String("eu-west-1"),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	return sess, err

}

//func CreateAwsSession2() (*session.Session, error) {
func CreateAwsSession2(c SlemConfig) (*ec2.EC2, error) {

	// Initialize a session in <region> that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: c.AWSProfile,

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(c.AWSRegion),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	svc := ec2.New(sess)

	return svc, err

}

var statusCMD = &cobra.Command{
	Use:   "status",
	Short: "Current status of setup / IP-allocaion",
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Printf("Do stuff later")

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)
		fmt.Printf("\nSTATUS environment (TAG GROUP: %s = %s)\n", ProjTagKey, ProjTagVal)
		PrintBar("VPC Info", "=", 120)
		ListVpc(svc)
		PrintBar("IGW Info", "-=-", 40)
		ListInternetGateway(svc)
		PrintBar("Subnet Info", "-=-", 40)
		ListSubnet(svc)
		PrintBar("Route Table Info", "-=-", 40)
		ListRTB(svc)
		PrintBar("Security Group Info", "-=-", 40)
		ListSG(svc)
		PrintBar("EC2 Instance Info", "===", 40)
		ListEC2(svc)

	},
}

func GetAllocatedV4(svc *ec2.EC2, ip4 string) []*ec2.Address {

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

	ret := result.Addresses
	return ret
}

type GroupConf struct {
	Domain string `yaml:"Domain"`
	IP4    string `yaml:"IP4"`
	IP6    string `yaml:"IP6"`
}

func LoadCourseConf(inFile string) (map[string]GroupConf, error) {

	yamlFile, err := ioutil.ReadFile(inFile)
	ret := map[string]GroupConf{}
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ret, nil
}

var validateCMD = &cobra.Command{
	Use:   "validate",
	Short: "Validate parameters i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		var inFile string
		if len(args) > 0 {
			inFile = args[0]
		} else {
			exitErrorf("No file specified\n")
		}

		errnr, _ := ValidateEnvironment(inFile)
		if errnr > 0 {
			fmt.Printf("Number of errors greater than 0 (%s), Cannot proceed\n", strconv.Itoa(errnr))
		}

	},
}

var deployCMD = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy environment described i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		SlemConf, _ := LoadSlemConf(awsconf)
		fmt.Printf("SlemConf = %+v\n", SlemConf)

		var tags []*ec2.Tag
		tags = ApendTags(tags, SlemConf.Tags)

		sess, err := CreateAwsSession()

		if err != nil {
			fmt.Printf("Session create error, %v", err)
		}

		// Create an EC2 service client.
		svc := ec2.New(sess)

		_ = DeployEnvironment(svc, tags)

	},
}

var destroyCMD = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy environment described i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Disassociating Addresses\n")
		fmt.Printf("Terminating EC2 Instances\n")
		fmt.Printf("Releasing Addresses\n")

		fmt.Printf("Removing Route from Routnig Table\n")
		fmt.Printf("Disassociating Routing Table from Subnet\n")
		fmt.Printf("Removing Routing Table\n")
		fmt.Printf("Removing Subnet\n")
		fmt.Printf("Detaching Internet Gateway from VPC\n")
		fmt.Printf("Removing Internet Gateway\n")
		fmt.Printf("Removing VPC\n")
		//check wait for instances to terminate
		fmt.Printf("Removing Security Group\n")

	},
}

func ValidateEnvironment(inFile string) (int, error) {
	sess, err := CreateAwsSession()

	if err != nil {
		fmt.Printf("Session create error, %v", err)
	}

	// Create an EC2 service client.
	svc := ec2.New(sess) //  type *ec2.EC

	gdata, _ := LoadCourseConf(inFile)
	errnr := 0
	for k, v := range gdata {
		if v.IP4 != "" {
			adr := GetAllocatedV4(svc, v.IP4)
			var msg, ok string
			if len(adr) < 1 {
				msg = fmt.Sprintf("(%s) IPv4 adress %s - available", k, v.IP4)
				ok = "(OK)"
			} else {
				for _, a := range adr {
					pip := *a.PublicIp
					msg = fmt.Sprintf("(%s) IPv4 adress %s - allocated", k, pip)
					if a.AssociationId != nil {
						asid := *a.AssociationId
						msg += fmt.Sprintf(" and associated (%s)", asid)
						ok = "(FAIL)"
						errnr++
					} else {
						msg += fmt.Sprintf(" but not associated")
						ok = "(OK)"
					}
				}
			}
			fmt.Printf("%-100s %s\n", msg, ok)

		}
	}
	return errnr, nil
}

// General functions?

func ValidIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func AddTag(svc *ec2.EC2, resourceId, tagName, tagValue string) error {
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{aws.String(resourceId)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tagName),
				Value: aws.String(tagValue),
			},
		},
	})
	if errtag != nil {
		exitErrorf("Could not create tags for instance\n %v", resourceId, errtag)
	}

	fmt.Printf("Successfully tagged resource %s\n", resourceId)
	return nil
}

func UnTag(svc *ec2.EC2, resourceId, tagName, tagValue string) error {
	_, errtag := svc.DeleteTags(&ec2.DeleteTagsInput{
		Resources: []*string{aws.String(resourceId)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tagName),
				Value: aws.String(tagValue),
			},
		},
	})
	if errtag != nil {
		exitErrorf("Could not create tags for instance\n %v", resourceId, errtag)
	}

	fmt.Printf("Successfully tagged resource %s", resourceId)
	return nil
}

func PrintBar(header, pattern string, length int) {
	if len(header) > 0 {
		header = pattern + "[ " + header + " ]"
		//patlen := len(pattern) * length
		headlen := len(header) / len(pattern)
		newlen := int(length - headlen)
		fmt.Printf("\n%s%s\n", header, strings.Repeat(pattern, newlen))
	} else {
		fmt.Printf("\n%s\n", strings.Repeat(pattern, length))
	}
}

func init() {

	rootCmd.AddCommand(awsCmd)
	awsCmd.AddCommand(
		ipcalcCMD, statusCMD, validateCMD, deployCMD, destroyCMD, addtagCMD, untagCMD,
		vpcCMD, igwCMD, sgCMD, rtbCMD, routeCMD, route6CMD, subnetCMD, ec2CMD, eipCMD, testCMD)

	// BUG: cannot update awsconf for some reason. Flag always at default.
	awsCmd.PersistentFlags().StringVarP(&awsconf, "awsconf", "", "./config/awsconf.yaml", "Basic configuration parameters for the SLEM AWS envirinment")
	deployCMD.Flags().StringVarP(&labconf, "conf", "c", "./config/labconf.yaml", "Parameters for the LAB Environments")
	addtagCMD.Flags().StringVarP(&resourceId, "resource-id", "r", "", "ID of resource to tag")
	addtagCMD.Flags().StringVarP(&tagName, "name", "", "", "Tag name (Key)")
	addtagCMD.Flags().StringVarP(&tagValue, "value", "", "", "Tag value (Value)")
	untagCMD.Flags().StringVarP(&resourceId, "resource-id", "r", "", "ID of resource to tag")
	untagCMD.Flags().StringVarP(&tagName, "name", "", "", "Tag name (Key)")
	untagCMD.Flags().StringVarP(&tagValue, "value", "", "", "Tag value (Value)")

	addtagCMD.MarkFlagRequired("resource-id")
	addtagCMD.MarkFlagRequired("name")
	addtagCMD.MarkFlagRequired("value")
	untagCMD.MarkFlagRequired("resource-id")
	untagCMD.MarkFlagRequired("name")
	untagCMD.MarkFlagRequired("value")

	ipcalcCMD.Flags().IntVarP(&ipstart, "start", "s", 1, "Start at this IP")
	ipcalcCMD.Flags().IntVarP(&iprange, "range", "r", 4, "Number of addresses per block")
	//ipcalc.Flags().StringVarP(&cidr, "cidr", "c", "/30", "Size of IP block")

	//groupCmd.PersistentFlags().StringVarP(&gin.day, "day", "D", "", "Day in YYYY-MM-DD format")

	SlemConf, _ := LoadSlemConf(awsconf)
	svc2, err = CreateAwsSession2(SlemConf)

}
