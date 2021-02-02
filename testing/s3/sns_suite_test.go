package s3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SNSSuiteTestSuite struct {
	SNSSuite
}

func TestSNSSuite(t *testing.T) {
	suite.Run(t, new(SNSSuiteTestSuite))
}

func (s *SNSSuiteTestSuite) TestHelpers() {
	fmt.Printf("+++++++++++ sns +++++++++++ \n")
	url := s.SNSSuite.CreateTopic("test")
	assert.NotEmpty(s.T(), url)
	s.SNSSuite.DeleteTopic(url)
}
