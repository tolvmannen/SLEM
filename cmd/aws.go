package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

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

	//cf       MainConfig
	e        Environment
	err      error
	xsIdErr  = errors.New("Multiple resource IDs returned")
	noIdErr  = errors.New("No resource IDs returned")
	ipAllErr = errors.New("IP address allocated")
	ipAssErr = errors.New("IP address associated")
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

		e.Load()
		e.MakeIpPlan()

		fmt.Printf("\n\n%+v\n", e)
		/*
			var tags []*ec2.Tag
			tags = ApendTags(tags, SlemConf.Tags)
		*/
	},
}

var deployCMD = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy environment described i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		e.Load()
		e.DeployEnv()
		e.MakeIpPlan()
		e.DeployEC2()

	},
}

var destroyCMD = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy environment described i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		e.Load()
		e.DestroyEC2()
		e.DestroyEnv()

	},
}

var statusCMD = &cobra.Command{
	Use:   "status",
	Short: "Current status of setup / IP-allocaion",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("\nSTATUS environment (Project = %s)\n", e.Project)
		PrintBar("VPC Info", "=", 120)
		e.ListVpc()
		PrintBar("IGW Info", "-=-", 40)
		e.ListIGW()
		PrintBar("Subnet Info", "-=-", 40)
		e.ListSubnet()
		PrintBar("Route Table Info", "-=-", 40)
		e.ListRTB()
		PrintBar("Security Group Info", "-=-", 40)
		e.ListSG()

		PrintBar("EC2 Instance Info", "===", 40)
		ListEC2(e.SVC)

	},
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

/*
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
*/

var validateCMD = &cobra.Command{
	Use:   "validate",
	Short: "Validate parameters i YAML file for course ",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("\nUpcoming feature, maybe..\n")

	},
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

	//cf, _ = LoadMainConf("./config/project-file.yaml")
	//cf.CreateSession()
	//cf.LoadEnvConf()

	e.LoadConf("./config/byoip.yaml")
	e.CreateSession()

}
