package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KinesisSuiteTestSuite struct {
	KinesisSuite
}

func TestKinesisSuite(t *testing.T) {
	suite.Run(t, new(KinesisSuiteTestSuite))
}

func (s *KinesisSuiteTestSuite) TestHelpers() {
	// Create a stream
	shardID := s.KinesisSuite.CreateStream("test")
	assert.NotEmpty(s.T(), shardID)

	recsInput := [][]byte{
		{1, 2},
		{3, 4},
	}
	s.KinesisSuite.PutRecords("test", recsInput)

	recs := s.KinesisSuite.GetRecords("test", shardID)
	assert.Equal(s.T(), recsInput, recs)

	s.KinesisSuite.DeleteStream("test")
}
