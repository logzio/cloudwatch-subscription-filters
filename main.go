package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	l "github.com/aws/aws-sdk-go/service/lambda"
	"go.uber.org/zap"
	lp "main/logger"
	"os"
	"strings"
)

var sugLog *zap.SugaredLogger

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event map[string]interface{}) (string, error) {
	logger := lp.GetLogger()
	defer logger.Sync()
	sugLog = logger.Sugar()
	sugLog.Info("Starting handling event...")
	sugLog.Debug("Handling event: ", event)
	err := validateRequired()
	if err != nil {
		return "Lambda finished with error", err
	}

	if _, ok := event["detail"]; ok {
		// Create log group invocation
		sugLog.Debug("Detected Eventbridge event")
		newLogGroupCreated(event["detail"].(map[string]interface{})["requestParameters"].(map[string]interface{})["logGroupName"].(string))
	} else {
		// First invocation
		if event["RequestType"].(string) == "Create" {
			sugLog.Debug("Detected Cloudformation initialization event")
			lambda.Start(cfn.LambdaWrap(customResourceRun))
		} else if event["RequestType"].(string) == "Delete" {
			sugLog.Debug("Detected Cloudformation delete event")
			lambda.Start(cfn.LambdaWrap(customResourceRunDelete))
		} else {
			lambda.Start(cfn.LambdaWrap(customResourceRunDoNothing))
		}
	}

	return "Lambda finished", nil
}

func validateRequired() error {
	destinationArn := os.Getenv(envShipperFuncArn)
	if destinationArn == emptyString {
		return fmt.Errorf("destination ARN must be set")
	}

	accountId := os.Getenv(envAccountId)
	if accountId == emptyString {
		return fmt.Errorf("account id must be set")
	}

	awsPartition := os.Getenv(envAwsPartition)
	if awsPartition == emptyString {
		return fmt.Errorf("aws partition must be set")
	}

	return nil
}

// Wrapper for first invocation from cloud formation custom resource
func customResourceRun(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	err = handleFirstInvocation()
	if err != nil {
		sugLog.Error("Error while handling first invocation: ", err.Error())
		return
	}

	return
}

// Wrapper for invocation from cloudformation custom resource - for read, update
func customResourceRunDoNothing(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}

// Wrapper for invocation from cloudformation custom resource - delete
func customResourceRunDelete(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	sess, err := getSession()
	if err != nil {
		sugLog.Error("Error while creating session: ", err.Error())
	}

	deleted := make([]string, 0)
	servicesToDelete := getServices()
	if servicesToDelete != nil {
		newDeleted, err := deleteServices(sess, servicesToDelete, deleted)
		deleted = append(deleted, newDeleted...)
		if err != nil {
			sugLog.Error(err.Error())
		}
	}

	pathsToDelete := getCustomPaths()
	if pathsToDelete != nil {
		newDeleted, err := deleteCustom(sess, pathsToDelete)
		deleted = append(deleted, newDeleted...)
		if err != nil {
			sugLog.Error(err.Error())
		}
	}

	sugLog.Debug("Deleted subscription filters for the following log groups: ", deleted)

	return
}

func newLogGroupCreated(logGroup string) {
	// Prevent a situation where we put subscription filter on the trigger and shipper function
	if logGroup == lambdaPrefix+getShipperFunctionName() ||
		logGroup == lambdaPrefix+os.Getenv(envFunctionName) {
		return
	}
	servicesToAdd := getServices()
	var added []string
	if servicesToAdd != nil {
		serviceToPrefix := getServicesMap()
		sess, err := getSession()
		if err != nil {
			sugLog.Error("Could not create aws session: ", err.Error())
			return
		}
		logsClient := cloudwatchlogs.New(sess)
		for _, service := range servicesToAdd {
			if prefix, ok := serviceToPrefix[service]; ok {
				if strings.Contains(logGroup, prefix) {
					added = putSubscriptionFilter([]string{logGroup}, logsClient)
					if len(added) > 0 {
						sugLog.Info("Added log group: ", logGroup)
						return
					}
				}
			}
		}
	}

	sugLog.Info("Log group ", logGroup, " does not match any of the selected services: ", servicesToAdd)
}

func handleFirstInvocation() error {
	sess, err := getSession()
	if err != nil {
		return err
	}

	err = addInvokePermissions(sess)
	if err != nil {
		return err
	}
	added := make([]string, 0)
	servicesToAdd := getServices()
	if servicesToAdd != nil {
		newAdded, err := addServices(sess, servicesToAdd)
		added = append(added, newAdded...)
		if err != nil {
			sugLog.Error(err.Error())
		}
	}

	pathsToAdd := getCustomPaths()
	if pathsToAdd != nil {
		newAdded, err := addCustom(sess, pathsToAdd, added)
		added = append(added, newAdded...)
		if err != nil {
			sugLog.Error(err.Error())
		}
	}

	sugLog.Debug("Following these log groups: ", added)

	return nil
}

func addInvokePermissions(sess *session.Session) error {
	principal := "logs.amazonaws.com"
	statementId := "log-groups-invoke"
	action := "lambda:InvokeFunction"
	accountId := os.Getenv(envAccountId)
	logGroupsArn := fmt.Sprintf("arn:%s:logs:%s:%s:log-group:*:*", os.Getenv(envAwsPartition), os.Getenv(envAwsRegion), accountId)
	shipperFuncName := getShipperFunctionName()
	lambdaClient := l.New(sess)
	addPermOutput, err := lambdaClient.AddPermission(&l.AddPermissionInput{
		Action:        &action,
		FunctionName:  &shipperFuncName,
		Principal:     &principal,
		SourceAccount: &accountId,
		SourceArn:     &logGroupsArn,
		StatementId:   &statementId,
	})

	if err != nil {
		if strings.Contains(err.Error(), "StatusCode: 409") {
			sugLog.Warn("Permissions already exist for the function.")
		} else {
			return fmt.Errorf("Error occurred while trying to add permission to shipper function: %s", err.Error())
		}
	} else {
		sugLog.Debug(addPermOutput.String())
	}

	return nil
}

func addCustom(sess *session.Session, customGroup, added []string) ([]string, error) {
	logsClient := cloudwatchlogs.New(sess)
	toAdd := make([]string, 0)
	lambdaNameTrigger := lambdaPrefix + os.Getenv(envFunctionName)
	lambdaNameShipper := lambdaPrefix + getShipperFunctionName()
	for _, customLogGroup := range customGroup {
		if !listContains(customLogGroup, added) {
			// Prevent a situation where we put subscription filter on the trigger and shipper function
			if customLogGroup != lambdaNameShipper &&
				customLogGroup != lambdaNameTrigger {
				toAdd = append(toAdd, customLogGroup)
			}
		}
	}

	newAdded := putSubscriptionFilter(toAdd, logsClient)

	return newAdded, nil
}

func addServices(sess *session.Session, servicesToAdd []string) ([]string, error) {
	logsClient := cloudwatchlogs.New(sess)
	logGroups := getLogGroups(servicesToAdd, logsClient)
	if len(logGroups) > 0 {
		sugLog.Debug("Detected the following services: ", logGroups)
		newAdded := putSubscriptionFilter(logGroups, logsClient)
		return newAdded, nil
	} else {
		return nil, fmt.Errorf("Could not retrieve any log groups")
	}
}

func putSubscriptionFilter(logGroups []string, logsClient *cloudwatchlogs.CloudWatchLogs) []string {
	destinationArn := os.Getenv(envShipperFuncArn)
	filterPattern := ""
	filterName := subscriptionFilterName
	added := make([]string, 0)
	for _, logGroup := range logGroups {
		_, err := logsClient.PutSubscriptionFilter(&cloudwatchlogs.PutSubscriptionFilterInput{
			DestinationArn: &destinationArn,
			FilterName:     &filterName,
			LogGroupName:   &logGroup,
			FilterPattern:  &filterPattern,
		})

		if err != nil {
			sugLog.Error("Error while trying to add subscription filter for ", logGroup, ": ", err.Error())
			continue
		}

		added = append(added, logGroup)
	}

	return added
}

func getLogGroups(services []string, logsClient *cloudwatchlogs.CloudWatchLogs) []string {
	logGroupsToAdd := make([]string, 0)
	serviceToPrefix := getServicesMap()
	for _, service := range services {
		if prefix, ok := serviceToPrefix[service]; ok {
			sugLog.Debug("Working on prefix: ", prefix)
			newLogGroups, err := logGroupsPagination(prefix, logsClient)
			if err != nil {
				sugLog.Errorf("Error while searching for log groups of %s: %s", service, err.Error())
			}

			logGroupsToAdd = append(logGroupsToAdd, newLogGroups...)
		} else {
			sugLog.Errorf("Service %s is not supported. Skipping.", service)
		}
	}

	return logGroupsToAdd
}

func logGroupsPagination(prefix string, logsClient *cloudwatchlogs.CloudWatchLogs) ([]string, error) {
	var nextToken *string
	logGroups := make([]string, 0)
	for {
		describeOutput, err := logsClient.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: &prefix,
			NextToken:          nextToken,
		})

		if err != nil {
			return nil, err
		}
		if describeOutput != nil {
			nextToken = describeOutput.NextToken
			for _, logGroup := range describeOutput.LogGroups {
				// Prevent a situation where we put subscription filter on the trigger and shipper function
				if *logGroup.LogGroupName == lambdaPrefix+os.Getenv(envFunctionName) ||
					*logGroup.LogGroupName == lambdaPrefix+getShipperFunctionName() {
					continue
				}

				logGroups = append(logGroups, *logGroup.LogGroupName)
			}
		}

		if nextToken == nil {
			break
		}
	}

	return logGroups, nil
}

func getSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(os.Getenv(envAwsRegion)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error occurred while trying to create a connection to aws: %s. Aborting", err.Error())
	}

	return sess, nil
}

func deleteServices(sess *session.Session, servicesToAdd, deleted []string) ([]string, error) {
	logsClient := cloudwatchlogs.New(sess)
	logGroups := getLogGroups(servicesToAdd, logsClient)
	if len(logGroups) > 0 {
		sugLog.Debug("Detected the following services: ", logGroups)
		newAdded := deleteSubscriptionFilter(logGroups, logsClient)
		return newAdded, nil
	} else {
		return nil, fmt.Errorf("Could not delete any log groups")
	}
}

func deleteSubscriptionFilter(logGroups []string, logsClient *cloudwatchlogs.CloudWatchLogs) []string {
	filterName := subscriptionFilterName
	deleted := make([]string, 0)
	for _, logGroup := range logGroups {
		_, err := logsClient.DeleteSubscriptionFilter(&cloudwatchlogs.DeleteSubscriptionFilterInput{
			FilterName:   &filterName,
			LogGroupName: &logGroup,
		})

		if err != nil {
			sugLog.Error("Error while trying to delete subscription filter for ", logGroup, ": ", err.Error())
			continue
		}

		deleted = append(deleted, logGroup)
	}

	return deleted
}

func deleteCustom(sess *session.Session, customGroup []string) ([]string, error) {
	logsClient := cloudwatchlogs.New(sess)
	newDeleted := deleteSubscriptionFilter(customGroup, logsClient)

	return newDeleted, nil
}
