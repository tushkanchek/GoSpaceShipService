package env

import (
	"os"
	"strconv"
)

type TestingConfig interface {
	FailProduceForTesting() bool
}

type testingConfig struct {
	failProduceForTesting bool
}

func (c *testingConfig) FailProduceForTesting() bool {
	return c.failProduceForTesting
}

func LoadTestingConfig() TestingConfig {
	failStr := os.Getenv("TEST_FAIL_PRODUCE")
	fail, _ := strconv.ParseBool(failStr)
	return &testingConfig{
		failProduceForTesting: fail,
	}
}
