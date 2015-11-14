package commands

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	ServicePath string `hcl:"service_path"`
	Rackspace   struct {
		User   string `hcl:"user"`
		Key    string `hcl:"key"`
		Region string `hcl:"region"`
	} `hcl:"rackspace"`
}

var DropshipCmd = &cobra.Command{
	Use:   "dropship",
	Short: "A tool for automated and distributed artifact deployment",
	Long: `Dropship allows servers to automatically check, download, and install
artifacts from a file repository in a distributed fashion.
	`,
}

var CfgFile string

func init() {
	DropshipCmd.PersistentFlags().StringVar(&CfgFile, "config", "/etc/dropship.d/dropship.hcl", "config file (default is path/config.yaml|json|toml)")
}

func Execute() {
	AddCommands()
	if err := DropshipCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func AddCommands() {
	DropshipCmd.AddCommand(agentCmd)
	DropshipCmd.AddCommand(versionCmd)
}

func InitializeConfig() {
	var cfg Config
	cfgData, err := ioutil.ReadFile(CfgFile)
	if err != nil {
		log.Fatalln("Unable to locate Config File. make sure you specify it using the --config flag")
		return
	}
	err = hcl.Decode(&cfg, string(cfgData))

	if err != nil {
		log.Fatalln("Unable to parse Config File.")
		return
	}

	viper.Set("servicePath", cfg.ServicePath)
	viper.Set("rackspaceUser", cfg.Rackspace.User)
	viper.Set("rackspaceKey", cfg.Rackspace.Key)
	viper.Set("rackspaceRegion", cfg.Rackspace.Region)
}
