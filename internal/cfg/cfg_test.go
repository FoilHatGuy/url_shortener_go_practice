package cfg

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestSequentialInitialisation() {
	config1 := Initialize()
	s.Assert().Equal(&config1, &config)
	config2 := Initialize()
	s.Assert().Equal(&config2, &config)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
