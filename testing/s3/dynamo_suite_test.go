package s3

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// test suite to test the helper functions in the Dynamo test suite base
// type
type DynamoSuiteTestSuite struct {
	DynamoSuite
}

// Right now this just makes sure that the dynamo container launches
// correctly. When helper functions are added to the suite they will
// be tested here.
func TestDynamoSuite(t *testing.T) {
	suite.Run(t, new(DynamoSuiteTestSuite))
}
