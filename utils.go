package main

import (
	"os"
	"strings"
)

const (
	envServices       = "SERVICES"
	envAwsRegion      = "AWS_REGION"               // reserved env
	envFunctionName   = "AWS_LAMBDA_FUNCTION_NAME" // reserved env
	envShipperFuncArn = "SHIPPER_ARN"
	envAccountId      = "ACCOUNT_ID"
	envCustomGroups   = "CUSTOM_GROUPS"
	envAwsPartition   = "AWS_PARTITION"

	valuesSeparator = ";"
	emptyString     = ""
	lambdaPrefix    = "/aws/lambda/"
)

func getServicesToAdd() []string {
	servicesStr := os.Getenv(envServices)
	if servicesStr == emptyString {
		return nil
	}

	servicesStr = strings.ReplaceAll(servicesStr, " ", "")
	return strings.Split(servicesStr, valuesSeparator)
}

func getServicesMap() map[string]string {
	return map[string]string{
		"apigateway":       "/aws/apigateway/",
		"rds":              "/aws/rds/cluster/",
		"cloudhsm":         "/aws/cloudhsm/",
		"cloudtrail":       "aws-cloudtrail-logs-",
		"codebuild":        "/aws/codebuild/",
		"connect":          "/aws/connect/",
		"elasticbeanstalk": "/aws/elasticbeanstalk/",
		"ecs":              "/aws/ecs/",
		"eks":              "/aws/eks/",
		"aws-glue":         "/aws/aws-glue/",
		"aws-iot":          "AWSIotLogsV2",
		"lambda":           "/aws/lambda/",
		"macie":            "/aws/macie/",
		"amazon-mq":        "/aws/amazonmq/broker/",
	}
}

func getCustomPaths() []string {
	pathsStr := os.Getenv(envCustomGroups)
	if pathsStr == emptyString {
		return nil
	}

	pathsStr = strings.ReplaceAll(pathsStr, " ", "")
	return strings.Split(pathsStr, valuesSeparator)
}

func listContains(s string, l []string) bool {
	for _, item := range l {
		if s == item {
			return true
		}
	}

	return false
}

func getShipperFunctionName() string {
	// ARN format is arn:aws:lambda:region:account-id:function:name
	arn := os.Getenv(envShipperFuncArn)
	arnArr := strings.Split(arn, ":")
	return arnArr[len(arnArr)-1]

}
