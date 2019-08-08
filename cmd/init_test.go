package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	env := map[string]string{
		"AzureSimple":   SIMPLE,
		"AzureCommon":   COMMON,
		"AzureSingleKV": KEYVAULT,
		"AzureMultiple": MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Testing Init function for environment", k)
		if error := Init(v, "test"+k); error != nil {
			t.Error("There was an error running the Init function.")
			return
		}
		if fileExists("bedrock/cluster/environments/test" + k + "/" + v + "/bedrock-config.tfvars") {
			fmt.Println("Configuration file successfully built for environment", k)
		} else {
			t.Error("There was an error with creating bedrock-config.tfvars")
		}
	}

	// Clean up test environment
	os.RemoveAll("bedrock")
}
