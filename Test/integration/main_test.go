//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestIntegration(t *testing.T) {
	// Get feature file from environment variable if specified
	featureFile := os.Getenv("INTEGRATION_FEATURE_FILE")
	// Get specific scenario tags from environment variable if specified
	scenarioTags := os.Getenv("INTEGRATION_SCENARIO_TAGS")

	var paths []string
	if featureFile != "" {
		// Run only the specified feature file
		paths = []string{"features/" + featureFile}
	} else {
		// Run all feature files
		paths = []string{"features"}
	}

	options := &godog.Options{
		Format:      "pretty",
		Concurrency: 1,
		Paths:       paths,
	}

	// Add tags filter if specific scenario tags are provided
	if scenarioTags != "" {
		options.Tags = scenarioTags
	}

	suite := godog.TestSuite{
		Name:                 "integration",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options:              options,
	}

	if exitCode := suite.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
