package instance

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stretchr/testify/assert"

	mocks "github.com/disneystreaming/go-ssmhelpers/testing"
)

func TestGetAllSSMInstances(t *testing.T) {
	assert := assert.New(t)

	// Set up our mock session and input object
	mockSvc := &mocks.MockSSMClient{}

	t.Run("latest agent version filter", func(t *testing.T) {
		ssmInput := &ssm.DescribeInstanceInformationInput{}
		instances, _ := GetAllSSMInstances(mockSvc, ssmInput, true)

		// Function should filter out any instances that are offline or not running Linux or do not have the latest agent version
		assert.Lenf(instances, 3, "Incorrect number of matching instances returned; got %d, expected 3", len(instances))
	})

	t.Run("standard instance filters", func(t *testing.T) {
		ssmInput := &ssm.DescribeInstanceInformationInput{}
		instances, _ := GetAllSSMInstances(mockSvc, ssmInput, false)

		// Function should filter out any instances that are offline or not running Linux
		assert.Lenf(instances, 4, "Incorrect number of matching instances returned; got %d, expected 4", len(instances))
	})

}
