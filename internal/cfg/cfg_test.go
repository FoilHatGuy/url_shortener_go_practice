package cfg

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestSequentialInitialization() {
	config1 := Initialize()
	s.Assert().Equal(&config1, &config)
	config2 := Initialize()
	s.Assert().Equal(&config2, &config)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
