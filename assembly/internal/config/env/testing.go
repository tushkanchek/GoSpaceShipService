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
	fail, err := strconv.ParseBool(failStr)
	if err != nil {
		// Если не установлена или некорректная, по умолчанию false
		fail = false
	}
	return &testingConfig{
		failProduceForTesting: fail,
	}
}
