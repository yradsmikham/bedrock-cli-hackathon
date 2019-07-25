package cmd

import (
	"errors"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"fmt"
	"math/rand"
	"time"
	. "github.com/logrusorgru/aurora"
)
var emoji_list = []string {
	":boom:", ":sparkles:", ":alien:", ":cat:", ":honeybee:", ":globe_with_meridians:", ":new_moon:", ":full_moon:", ":earth_americas:", ":earth_asia:", ":tropical_fish:", ":penguin:", ":baby_chick:", ":koala:", ":zap:", ":cyclone:", ":dog:", ":bear:", ":panda_face:", ":maple_leaf:", ":mushroom:", ":full_moon_with_face:", ":crescent_moon:", ":snowflake:", ":frog:", ":monkey_face:", ":snail:", ":rabbit2:", ":new_moon_with_face:", ":bulb:", ":floppy_disk:", ":tennis:", ":gem:", ":baby_bottle:", ":birthday:", ":green_apple:", ":basketball:", ":coffee:", ":tangerine:", ":soccer:", ":game_die:", ":tea:", ":cookie:", ":tomato:", ":lemon:", ":pizza:", ":apple:", ":doughnut:", ":package:", ":dvd:", ":baseball:",":dart:",
}
var info_map = map[string]map[string][]string {
	SIMPLE: {
		"info": []string{
			Bold(Green(SIMPLE)).String() + " environment is a non-production ready template provided to easily try out Bedrock on Azure",
			"Deploys a single cluster (with Flux) using a service principal of your choice",
		},
		"pre-reqs": []string{
			"Service Principal: You can generate an azure service principal using the " + Bold(Green("az ad sp create-for-rbac --subscription <id | name>")).String() + " command",
			"A Kubernetes manifest repository",
		},
	},
	MULTIPLE: {
		"info": []string{
			Bold(Green(MULTIPLE)).String() + " environment deploys three redundant clusters (with Flux on each cluster) and an Azure Keyvault, each behind Azure Traffic Manager, which is configured with rules for routing traffic to one of the three clusters",
			"The Public IP for each AKS cluster will be provisioned in the Resource Group for each region",
			"A Traffic Manager Rule will be created for each Public IP Address so that the Traffic Manager knows about and can route traffic accordingly",
			"By default, the multiple cluster template has configurations set up for aks-eastus, aks-westus and aks-centralus. If your regional requirements differ, modify these names to match",
			"Each cluster uses its own resource group and resource group location",
			"Each cluster uses its own gitops path (although each cluster can still point to the same path)",
		},
		"pre-reqs": []string{
			"Dependent on a successful deploment of " + Bold(Green(COMMON)).String(),
			"Service Principal needs to have Owner privileges on the Azure subscription",
			"Traffic Manager's following properties are required: Profile name, DNS name, resource group name and resource group location",
			"A Kubernetes manifest repository",
		},
	},
	KEYVAULT: {
		"info": []string{
			Bold(Green(KEYVAULT)).String() + " environment deploys a single production level AKS cluster configured with Flux and Azure Keyvault",
		},
		"pre-reqs": []string{
			"Dependent on a successful deploment of " + Bold(Green(COMMON)).String(),
			"A Kubernetes manifest repository",
		},
	},
	COMMON: {
		"info": []string{
			Bold(Green(COMMON)).String() + " environment is a production ready template to setup common permanent elements of your infrastructure like vnets, keyvault, and a common resource group for them",
			"Dependency environment for other environments like the " + KEYVAULT + " and " + MULTIPLE,
			"Creates a resource group for your deployment, a VNET and subnet(s), and an Azure Key Vault with the appropriate access policies",
		},
		"pre-reqs": []string{
			"A storage account in Azure: set the following fields as environment variables or pass as parameters: " + Bold(Green("AZURE_STORAGE_ACCOUNT")).String() + ", " + Bold(Green("AZURE_STORAGE_KEY")).String() + ", " + Bold(Green("ARM_SUBSCRIPTION_ID")).String() + ", " + Bold(Green("ARM_CLIENT_ID")).String() + ", " + Bold(Green("ARM_CLIENT_SECRET")).String() + ", " + Bold(Green("ARM_TENANT_ID")).String(),
		},
	},
}

func GetEmoji() (emoji string) {
	rand.Seed(time.Now().UnixNano())
	rand_num := rand.Intn(len(emoji_list))
	return emoji_list[rand_num]
}

func Info(env string) (err error) {
	var emoji_str = GetEmoji()

	fmt.Println();
	for _, element := range info_map[env]["info"] {
		fmt.Println(emoji.Sprintf("%s %s", emoji_str, element))
	}
	fmt.Println(Bold(Cyan("\n    Pre-Requisites")))
	for _, element := range info_map[env]["pre-reqs"] {
		fmt.Println(emoji.Sprintf("%s %s", emoji_str, element))
	}
	return nil
}

var infoCmd = &cobra.Command{
	Use:   "info <environment_name>",
	Short: "Get details about an environment",
	Long:  `Get details about an environment`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if (len(args) == 0) {
			return errors.New("You need to specify an environment: " + SIMPLE + ", " + MULTIPLE + ", " + KEYVAULT + ", " + COMMON)

		}
		if !((args[0] == SIMPLE) || (args[0] == MULTIPLE) || args[0] == KEYVAULT || args[0] == COMMON) {
			return errors.New("The environment you specified is not of the following: " + SIMPLE + ", " + MULTIPLE + ", " + KEYVAULT + ", " + COMMON)
		}
		return Info(args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
