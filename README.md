# cloudwatch-shipper

AWS Lambda function that ships Cloudwatch logs to logz.io

## Instructions

To deploy this project, click the button that matches the region you wish to deploy your Stack to:

| Region           | Deployment                                                                                                                                                                                                                                                                                                                                             |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `us-east-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-east-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `us-east-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-east-2.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `us-west-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-west-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-west-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `us-west-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-us-west-2.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `eu-central-1`   | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-central-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-central-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)     | 
| `eu-north-1`     | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-north-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-north-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)         | 
| `eu-west-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `eu-west-2`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-2.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `eu-west-3`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=eu-west-3#/stacks/create/review?templateURL=https://logzio-aws-integrations-eu-west-3.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `sa-east-1`      | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=sa-east-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-sa-east-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)           | 
| `ap-northeast-1` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper) | 
| `ap-northeast-2` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-2.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper) | 
| `ap-northeast-3` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-northeast-3#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-northeast-3.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper) | 
| `ap-south-1`     | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-south-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-south-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)         | 
| `ap-southeast-1` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-southeast-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper) | 
| `ap-southeast-2` | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-2#/stacks/create/review?templateURL=https://logzio-aws-integrations-ap-southeast-2.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper) | 
| `ca-central-1`   | [![Deploy to AWS](https://dytvr9ot2sszz.cloudfront.net/logz-docs/lights/LightS-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ca-central-1#/stacks/create/review?templateURL=https://logzio-aws-integrations-ca-central-1.s3.amazonaws.com/cloudwatch-shipper/0.0.1/sam-template.yaml&stackName=logzio-cloudwatch-shipper)     |

### 1. Specify stack details

Specify the stack details as per the table below, check the checkboxes and select **Create stack**.

| Parameter          | Description                                                                                                                                       | Required/Default |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|------------------|
| `logGroupName`     | The name of log group you want to ship logs from to Logz.io                                                                                       | **Required**     |
| `logzioToken`      | Your Logz.io log shipping token.                                                                                                                  | **Required**     |
| `logzioListener`   | The Logz.io listener URL fot your region. (For more details, see the [regions page](https://docs.logz.io/user-guide/accounts/account-region.html) | **Required**     |
| `logType`          | The log type you'll use with this Lambda. This is shown in your logs under the type field in Kibana. Logz.io applies parsing based on type.       | `cloudwatch`     |
| `logLevel`         | Log level for the Lambda function. Can be one of: debug, info, warn, error, fatal, panic.                                                         | `info`           |
| `compress`         | If true, the Lambda will send compressed logs. If false, the Lambda will send uncompressed logs.                                                  | `true`           |
| `sendAll`          | By default, we do not send logs of type START, END, REPORT. Choose true to send all log types.                                                    | `false`          |
| `additionalFields` | Enriches the CloudWatch events with custom properties at ship time. The format is `key1=value1;key2=value2`. By default is empty.                 | -                |

### 2. Send logs

Give the stack a few minutes to be deployed.

Once new logs are added to your chosen log group, they will be sent to your Logz.io account.


## Changelog

- **0.0.1**: Initial release.