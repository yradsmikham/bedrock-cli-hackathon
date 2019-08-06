package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	env := map[string]string{
		"AzureSimple":           SIMPLE,
		"AzureCommonInfra":      COMMON,
		"AzureSingleKeyVault":   KEYVAULT,
		"AzureMultipleClusters": MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Testing Init function for environment", k)
		Init(v, "test"+k)
		if fileExists("bedrock/cluster/environments/test" + k + "/" + v + "/bedrock-config.tfvars") {
			fmt.Println("Configuration file successfully built for environment", k)
		} else {
			t.Error("There was an error with creating bedrock-config.tfvars")
		}
	}

	// Clean up test environment
	os.RemoveAll("bedrock")
}
