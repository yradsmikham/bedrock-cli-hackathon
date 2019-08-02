package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yradsmikham/bedrock-cli/utils"
)

func setEnv(name string) (err error) {
	// must environment variables from bedrock-config and set them as environment variables
	viper.SetConfigName("bedrock-config")             // name of config file (without extension)
	viper.AddConfigPath(name + "/azure-common-infra") // path to look for the config file in
	errr := viper.ReadInConfig()                      // Find and read the config file
	if errr != nil {                                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", errr))
	}

	log.Info(emoji.Sprintf(":arrows_clockwise: Setting Environments Variables..."))
	os.Setenv("ARM_SUBSCRIPTION_ID", viper.GetString("subscription"))
	os.Setenv("ARM_CLIENT_ID", viper.GetString("service_principal_id"))
	os.Setenv("ARM_CLIENT_SECRET", viper.GetString("secret"))
	os.Setenv("ARM_TENANT_ID", viper.GetString("tenant_id"))

	return err
}

// Simulate or dry-run a bedrock environment creation (azure simple, multi-cluster, keyvault, etc.)
func Simulate(name string) (err error) {
	log.Info(emoji.Sprintf(":beginner: Starting Environment Deployment Simulation!"))

	files, err := ioutil.ReadDir(name)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		log.Info(emoji.Sprintf(":eyes: Searching for Azure-Common-Infra environment..."))
		if f.Name() == "azure-common-infra" {
			log.Info(emoji.Sprintf(":round_pushpin: Azure-Common-Infra environment found!"))
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Common-Infra Environment"))
			setEnv(name)

			// Terraform Init
			utils.TerraformInitBackend(name + "/azure-common-infra")

			// Terraform Plan
			utils.TerraformPlan(name + "/azure-common-infra")

			break
		}
	}

	// Run Terraform Init on everything else (e.g. azure-single-keyvault, azure-multi-cluster)
	for _, f := range files {
		if f.Name() == "azure-simple" {
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Simple Environment"))

			// Terraform Init
			utils.TerraformInit(name + "/azure-simple")

			// Terraform Plan
			utils.TerraformPlan(name + "/azure-simple")

			break
		}
		if f.Name() == "azure-single-keyvault" {
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Single-Keyvault Environment"))
			setEnv(name)

			// Terraform Init
			utils.TerraformInitBackend(name + "/azure-single-keyvault")

			// Terraform Plan
			utils.TerraformPlan(name + "/azure-single-keyvault")

			break
		}
		if f.Name() == "azure-multiple-clusters" {
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Multiple-Clusters Environment"))
			setEnv(name)

			// Terraform Init
			utils.TerraformInitBackend(name + "/azure-multiple-clusters")

			// Terraform Plan
			utils.TerraformPlan(name + "/azure-multiple-clusters")

			break
		}
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Completed simulated dry-run of environment deployment!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock deploy " + name))
	}

	return err
}

var simulateCmd = &cobra.Command{
	Use:   "simulate <environment-name>",
	Short: "Simulate the environment deployment using Terraform",
	Long:  `Simulate the environment deployment using terraform init and plan`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var name = "unique-environment-name"

		if len(args) > 0 {
			name = args[0]
		}
		return Simulate(name)
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}
