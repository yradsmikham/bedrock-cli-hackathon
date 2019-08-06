package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestSimulate(t *testing.T) {
	// Runs Init Function
	env := map[string]string{
		"az-simple":       SIMPLE,
		"az-common-infra": COMMON,
		"az-single-kv":    KEYVAULT,
		"az-multiple":     MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Testing Init function for environment", k)
		_, err := Init(v, "test-"+k)
		if err != nil {
			t.Error("There was an error creating test environment", k)
		}
		errSim := Simulate("bedrock/cluster/environments/test-" + k)
		if errSim != nil {
			t.Error("There was an error simulating", k)
		}
	}

	// Clean up test environment
	os.RemoveAll("bedrock")
}