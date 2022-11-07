package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose, debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slem",
	Short: "Simple Lab Environment Manager",
	Long:  `A handy tool to set up, take down and manage your lab environment.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("slem called\n")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "awsconf", "", "Config file for AWS account and basic parameters(default is ./AWS/awsconf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug output")
}

func initConfig() {
	// Mayhaps add things here later...

}
