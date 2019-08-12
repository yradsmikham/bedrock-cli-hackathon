package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestSimulate(t *testing.T) {
	// Runs Init Function
	env := map[string]string{
		// names in lowercase because azure-multiple-clusters env requires a domain_name_label
		//that uses only lowercase alphanumeric characters, numbers and hyphens
		"azuresimple":   SIMPLE,
		"azurecommon":   COMMON,
		"azuresinglekv": KEYVAULT,
		"azuremultiple": MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Test simulation for environment", k)
		errInit := Init(v, "test-"+k)
		if errInit != nil {
			t.Error("There was an error creating test environment", k)
		}
		errSim := Simulate("bedrock/cluster/environments/test-" + k)
		if errSim != nil {
			t.Error("There was an error simulating", k)
			t.Error(errSim)
		}
	}

	// Clean up test environment
	os.RemoveAll("bedrock")
}
