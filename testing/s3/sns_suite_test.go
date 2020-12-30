package s3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SNSSuite_TestSuite struct {
	SNSSuite
}

func TestSNSSuite(t *testing.T) {
	suite.Run(t, new(SNSSuite_TestSuite))
}

func (s *SNSSuite_TestSuite) TestHelpers() {
	fmt.Printf("+++++++++++ sns +++++++++++ \n")
	url := s.SNSSuite.CreateTopic("test")
	assert.NotEmpty(s.T(), url)
	s.SNSSuite.DeleteTopic(url)
}
