package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestSimulate(t *testing.T) {
	// Runs Init Function
	env := map[string]string{
		"AzureSimple":   SIMPLE,
		"AzureCommon":   COMMON,
		"AzureSingleKV": KEYVAULT,
		"AzureMultiple": MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Test simulation for environment", k)
		err := Init(v, "test"+k)
		if err != nil {
			t.Error("There was an error creating test environment", k)
		}
		errSim := Simulate("bedrock/cluster/environments/test" + k)
		if errSim != nil {
			t.Error("There was an error simulating", k)
			t.Error(errSim)
		}
	}

	// Clean up test environment
	os.RemoveAll("bedrock")
}
