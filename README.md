# Cloudwatch shipper with log groups detection

This project deploys instrumentation that allows shipping Cloudwatch logs to Logz.io.

## Overview:

This project will create 2 Lambda functions:

- **Shipper function**: this function is responsible for processing and shipping the Cloudwatch logs to Logz.io. [See here the function's repo](https://github.com/logzio/logzio_aws_serverless/tree/master/python3/cloudwatch).

- **Trigger function**: this function is responsible for adding subscription filters to the desired Cloudwatch log groups, to trigger the shipper function.

When the Trigger function is run for the first time, it will add subscription filters to the log groups chosen by the user. If the user chose a service, the Trigger function will also get triggered whenever a log group is created to check if this log group is for a service that is one of the services that the user has selected. If yes, it will add a subscription filter to it.

## Instructions

To deploy this project, click the button that matches the region you wish to deploy your Stack to:

| Region           | Deployment                                                                                                                                                                                                                                                                                                                                             |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `us-east-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-east-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `us-east-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-east-2.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `us-west-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-west-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-west-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `us-west-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-west-2.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `eu-central-1`   | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-central-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-central-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)     | 
| `eu-north-1`     | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-north-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-north-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)         | 
| `eu-west-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `eu-west-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-2.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `eu-west-3`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-3#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-3.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `sa-east-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=sa-east-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-sa-east-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)           | 
| `ap-northeast-1` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch) | 
| `ap-northeast-2` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-2.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch) | 
| `ap-northeast-3` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-3#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-3.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch) | 
| `ap-south-1`     | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-south-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-south-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)         | 
| `ap-southeast-1` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-southeast-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch) | 
| `ap-southeast-2` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-southeast-2.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch) | 
| `ca-central-1`   | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ca-central-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ca-central-1.s3.amazonaws.com/cloudwatch-shipper-trigger/1.1.1/sam-template.yaml&stackName=logzio-cloudwatch)     |


### 1. Specify stack details

Specify the stack details as per the table below, check the checkboxes and select **Create stack**.

#### Shipper config:

| Parameter              | Description                                                                                                                                                                                                                             | Required/Default                  |
|------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------|
| `logzioToken`          | Replace `<<SHIPPING-TOKEN>>` with the [token](https://app.logz.io/#/dashboard/settings/general) of the account you want to ship to.                                                                                                     | **Required**                      |
| `logzioListener`       | Listener host, and port (for example, `https://<<LISTENER-HOST>>:8071`).                                                                                                                                                                | **Required**                      |
| `logzioType`           | The log type you'll use with this Lambda. This can be a [built-in log type](https://docs.logz.io/user-guide/log-shipping/built-in-log-types.html), or a custom log type. <br> You should create a new Lambda for each log type you use. | Default: `logzio_cloudwatch_logs` |
| `logzioFormat`         | `json` or `text`. If `json`, the Lambda function will attempt to parse the message field as JSON and populate the event data with the parsed fields.                                                                                    | Default: `text`                   |
| `logzioCompress`       | Set to `true` to compress logs before sending them. Set to `false` to send uncompressed logs.                                                                                                                                           | Default: `true`                   |
| `logzioEnrich`         | Enrich CloudWatch events with custom properties, formatted as `key1=value1;key2=value2`.                                                                                                                                                | -                                 |
| `shipperLambdaTimeout` | The number of seconds that Lambda allows a function to run before stopping it, for the shipper function.                                                                                                                                | Default: `60`                     |
| `shipperLambdaMemory`  | Shipper function's allocated CPU proportional to the memory configured, in MB.                                                                                                                                                          | `512`                             |
| `shipperLogLevel` (Default: `INFO`)              | Log level for the shipper function. Possible values are: `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`.                                                                                                                               |
| `shipperRequestTimeout`  (Default: `15`)                | Timeout in seconds for each http request for sending logs into logz.io.                                                                                                                                                                 |


#### Trigger config:

| Parameter               | Description                                                                                                                                                                                                                                              | Required/Default |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------|
| `services`              | A comma-seperated list of services you want to collect logs from. Supported options are: `apigateway`, `rds`, `cloudhsm`, `cloudtrail`, `codebuild`, `connect`, `elasticbeanstalk`, `ecs`, `eks`, `aws-glue`, `aws-iot`, `lambda`, `macie`, `amazon-mq`. | -                |
| `customLogGroups`       | A comma-seperated list of custom log groups you want to collect logs from                                                                                                                                                                                | -                |
| `triggerLambdaTimeout`  | The amount of seconds that Lambda allows a function to run before stopping it, for the trigger function.                                                                                                                                                 | `60`             |
| `triggerLambdaMemory`   | Trigger function's allocated CPU proportional to the memory configured, in MB.                                                                                                                                                                           | `512`            |
| `triggerLambdaLogLevel` | Log level for the Lambda function. Can be one of: `debug`, `info`, `warn`, `error`, `fatal`, `panic`                                                                                                                                                     | `info`           |

##### ⚠️ Important note ⚠️

AWS limits every log group to have up to 2 subscription filters. If your chosen log group already has 2 subscription filters, the trigger function won't be able to add another one.

### 2. Send logs

Give the stack a few minutes to be deployed.

Once new logs are added to your chosen log group, they will be sent to your Logz.io account.

##### ⚠️ Important note ⚠️

If you've used the `services` field, you'll have to **wait 6 minutes** before creating new log groups for your chosen services. This is due to cold start and custom resource invocation, that can cause the cause Lambda to behave unexpectedly.


### Changelog:
- **1.1.1**:
  - Upgrade to `cloudwatch-shipper` 1.1.1:
    - Support Lambda insights.
    - Add configurable request timeout for shipper function.
    - Support configuring log level for shipper function.
- **1.1.0**:
  - Upgrade to `cloudwatch-shipper` 1.0.0:
  - **Breaking changes**:
    - For auto-detection of log level - log level will appear in upper case.
  - Lambda logs - send all logs include platform logs (START, END, REPORT). 
  - Add namespace field to logs - service name based on the log group name. 
- **1.0.0**: Initial release.
