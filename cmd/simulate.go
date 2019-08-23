package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	util "github.com/yradsmikham/bedrock-cli/util"
)

func setEnv(name string, env string) {
	// must retreive environment variables from bedrock-config and set them as environment variables
	viper.SetConfigName("bedrock-sp-config") // name of config file (without extension)
	viper.AddConfigPath(name + "/" + env)    // path to look for the config file in

	errr := viper.ReadInConfig() // Find and read the config file
	if errr != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", errr))
	}
	log.Info(emoji.Sprintf(":arrows_clockwise: Setting Environments Variables..."))
	os.Setenv("ARM_SUBSCRIPTION_ID", viper.GetString("subscription"))
	os.Setenv("ARM_CLIENT_ID", viper.GetString("service_principal"))
	os.Setenv("ARM_CLIENT_SECRET", viper.GetString("secret"))
	os.Setenv("ARM_TENANT_ID", viper.GetString("tenant_id"))
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
		if f.Name() == COMMON {
			log.Info(emoji.Sprintf(":round_pushpin: Azure-Common-Infra environment found!"))
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Common-Infra Environment"))
			setEnv(name, COMMON)

			// Terraform Init
			if error := util.TerraformInitBackend(name + "/azure-common-infra"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformPlan(name + "/azure-common-infra"); error != nil {
				return error
			}

			break
		}
	}

	// Run Terraform Init on everything else (e.g. azure-single-keyvault, azure-multi-cluster)
	for _, f := range files {
		if f.Name() == SIMPLE {
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Simple Environment"))
			setEnv(name, SIMPLE)

			// Terraform Init
			if error := util.TerraformInit(name + "/azure-simple"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformPlan(name + "/azure-simple"); error != nil {
				return error
			}

			break
		}
		if f.Name() == KEYVAULT {
			log.Info(emoji.Sprintf(":rocket: Deploying Azure Common Infra environment"))

			// Need to deploy azure-common-infra before you can run `terraform init` or `terraform apply`
			if error := util.TerraformApply(name + "/azure-common-infra"); error != nil {
				return error
			}
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Single-Keyvault Environment"))
			setEnv(name, KEYVAULT)
			// Terraform Init
			if error := util.TerraformInitBackend(name + "/azure-single-keyvault"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformPlan(name + "/azure-single-keyvault"); error != nil {
				return error
			}

			break
		}
		if f.Name() == MULTIPLE {
			log.Info(emoji.Sprintf(":dancers: Simulating Azure-Multiple-Clusters Environment"))
			setEnv(name, MULTIPLE)

			// Terraform Init
			if error := util.TerraformInit(name + "/azure-multiple-clusters"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformPlan(name + "/azure-multiple-clusters"); error != nil {
				return error
			}

			break
		}
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Completed simulated dry-run of environment deployment!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock deploy " + name + "'"))
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
