package worktime

import (
	"testing"

	"github.com/Vehsamrak/worktime/src/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestApplication(test *testing.T) {
	suite.Run(test, &ApplicationTestSuite{})
}

type ApplicationTestSuite struct {
	suite.Suite
}

func (suite *ApplicationTestSuite) Test_getHelpMessage_noParameters_helpMessageReturned() {
	application := &Application{}

	helpMessage := application.getHelpMessage()

	_, ok := helpMessage.(message.Message)
	assert.True(suite.T(), ok)
}
