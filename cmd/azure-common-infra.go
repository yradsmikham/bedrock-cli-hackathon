package cmd

import (
	"github.com/spf13/cobra"
)

var tenant string

// Initializes the configuration for the given environment
func commonInfra(servicePrincipal string, tenant string) (err error) {
	Init(COMMON)
	return err
}

var commonInfraCmd = &cobra.Command{
	Use:   COMMON + " --sp service-principal-app-id --tenant tenant-id",
	Short: "Deploys the Bedrock Common Infra Environment",
	Long:  `Deploys the Bedrock Common Infra Environment`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return commonInfra(servicePrincipal, tenant)
	},
}

func init() {
	commonInfraCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	commonInfraCmd.Flags().StringVar(&tenant, "tenant", "", "Password for the Service Principal")
	commonInfraCmd.MarkFlagRequired("sp")
	commonInfraCmd.MarkFlagRequired("tenant")
	rootCmd.AddCommand(commonInfraCmd)
}
