//go:build integration
// +build integration

package integration

import (
	"github.com/cucumber/godog"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "integration",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options: &godog.Options{
			Format:      "pretty",
			Concurrency: 1,
			Paths:       []string{"features"},
		},
	}

	if exitCode := suite.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
