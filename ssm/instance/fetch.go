package instance

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/disneystreaming/go-ssmhelpers/util/batch"
)

// getSSMInstances takes an SSM session and an *ssm.DescribeInstanceInformationInput object to get information about a given SSM-managed EC2 instance.
func getSSMInstances(session ssmiface.SSMAPI, input *ssm.DescribeInstanceInformationInput) (output *ssm.DescribeInstanceInformationOutput, err error) {
	return session.DescribeInstanceInformation(input)
}

// getSSMInstances takes an SSM session and an *ssm.DescribeInstanceInformationInput object to get information about a given SSM-managed EC2 instance.
func getSSMInstancesPages(session ssmiface.SSMAPI, input *ssm.DescribeInstanceInformationInput) (output []*ssm.DescribeInstanceInformationOutput, err error) {
	results, err := getSSMInstances(session, input)
	if err != nil {
		return nil, err
	}

	output = append(output, results)

	for results.NextToken != nil && *results.NextToken != "" && err == nil {
		input.SetNextToken(*results.NextToken)
		results, err = getSSMInstances(session, input)
		if err != nil {
			return output, err
		}
		output = append(output, results)
	}

	return output, err
}

// GetAllSSMInstances queries the SSM API for information about SSM-managed instances, and filters out any instances that are unresponsive to ping or not running Linux.
// It can also check to determine if instances are running the latest SSM agent version, which is a prerequisite for connecting via start-session functionality.
func GetAllSSMInstances(session ssmiface.SSMAPI, input *ssm.DescribeInstanceInformationInput, checkLatestAgent bool) (output []*ssm.InstanceInformation, err error) {
	results, err := getSSMInstancesPages(session, input)

	instanceIds := []*string{}
	for _, page := range results {
		for _, instanceInfo := range page.InstanceInformationList {
			instanceIds = append(instanceIds, instanceInfo.InstanceId)
		}
	}

	// Create our batch function to iterate over the first batch of results
	getInstanceBatch := func(min int, max int) (bool, error) {
		// Set up our additional filters and second input
		statusInput := &ssm.DescribeInstanceInformationInput{
			Filters: []*ssm.InstanceInformationStringFilter{
				{
					Key:    aws.String("PingStatus"),
					Values: aws.StringSlice([]string{"Online"}),
				},
				{
					Key:    aws.String("PlatformTypes"),
					Values: aws.StringSlice([]string{"Linux"}),
				},
				{
					Key:    aws.String("InstanceIds"),
					Values: instanceIds[min:max],
				},
			},
		}

		// Re-run query to filter out instances that are unresponsive to ping or not running Linux
		results, err = getSSMInstancesPages(session, statusInput)
		if err != nil {
			return false, err
		}

		for _, instanceList := range results {
			for _, v := range instanceList.InstanceInformationList {
				if checkLatestAgent && *v.IsLatestVersion {
					output = append(output, v)
				} else if !checkLatestAgent {
					output = append(output, v)
				}
			}
		}
		return true, nil
	}

	err = batch.Chunk(len(instanceIds), 100, getInstanceBatch)

	return output, err
}
