# aws-servicemap
Not all AWS services are supported in all regions. 

Use this go module to determine if a service is supported in a specific region or not.  

This module uses https://api.regional-table.region-services.aws.a2z.com/index.json as the source of truth.

## Functions

* [GetRegionsForService](#GetRegionsForService)
* [GetServicesForRegion](#GetServicesForRegion)
* [GetAllRegions](#GetAllRegions)
* [GetAllServices](#GetAllServices)
* [IsServiceInRegion](#IsServiceInRegion)

## Examples
### GetRegionsForService

```
package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	regions := awsservicemap.GetRegionsForService("grafana")
	fmt.Println(regions)
}
```

Output: 
```
[ap-northeast-1 ap-northeast-2 ap-southeast-1 ap-southeast-2 eu-central-1 eu-west-1 eu-west-2 us-east-1 us-east-2 us-west-2]
```


### GetServicesForRegion

```
package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	services := awsservicemap.GetServicesForRegion("af-south-2")
	fmt.Println(services)
}
```

Output: 
```
[route53 ssm codedeploy ec2 cloudtrail redshift iam es sqs apigateway stepfunctions privatelink acm artifact cloudformation kinesis phd rds trustedadvisor vpc directconnect aurora ebs fargate eventbridge xray emr secretsmanager swf marketplace cloudwatch dms organizations elasticache kms support sns vpn autoscaling ecr elb s3 dynamodb config cloudfront lambda ecs]
```

### GetAllRegions

```
package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	totalRegions := awsservicemap.GetAllRegions()
	fmt.Println(totalRegions)
}
```

Output: 
```
[ap-east-1 ap-northeast-1 ap-northeast-2 ap-south-1 ap-southeast-2 ca-central-1 eu-west-2 eu-west-3 us-east-1 us-east-2 ap-southeast-1 eu-central-1 eu-north-1 eu-west-1 us-gov-west-1 us-west-1 us-west-2 ap-northeast-3 cn-north-1 eu-central-2 sa-east-1 af-south-1 ap-southeast-3 eu-south-1 eu-south-2 me-central-1 cn-northwest-1 me-south-1 us-gov-east-1]
```

### GetAllServices

```
package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	totalServices := awsservicemap.GetAllServices()
	fmt.Println(totalServices)
}
```

Output: 
```
[translate route53 lookoutmetrics opsworkspuppetenterprise datasync frauddetector shield mwaa cloudshell servicediscovery sagemaker elastictranscoder devops-guru transitgateway ivs rosa elastic-inference ssm codedeploy iotdevicedefender ahl ec2 inspector2 quicksight cloudtrail greengrass polly athena rdsvmware groundstation chatbot redshift iam iot1click-projects filecache iotevents kinesisanalytics honeycode es dataexchange guardduty nimble vmwarecloudonaws firehose kinesisvideo appstream deepcomposer mq ram cloudsearch sso managedservices iotanalytics wellarchitectedtool sqs compute-optimizer braket waf glue outposts medialive apigateway proton stepfunctions budgets license-manager ds privatelink acm personalize artifact eks fms workspaces-web cloudformation kinesis cloudenduredisasterrecovery clouddirectory cognito-identity phd datapipeline controltower lumberyard rds grafana wam codeartifact fsx-ontap detective lightsail iotsitewise kafka trustedadvisor vpc directconnect aurora iot ebs chime aiq cloudenduremigration resiliencehub comprehend fargate backup mediapackage globalaccelerator snowcone drs kendra devicefarm eventbridge lex-runtime appflow xray textract cloudhsmv2 neptune fis amplify auditmanager emr workdocs secretsmanager swf augmentedairuntime marketplace batch mgn transfer ses codepipeline application-autoscaling timestream lakeformation mgh transcribemedical opsworks cloudwatch opsworkschefautomate amazonlocationservice codestar cloud9 workspaces managedblockchain snowball dms serverlessrepo robomaker pinpoint elasticbeanstalk transcribe organizations iotdevicemanagement macie aps elasticache mcs connect kms forecast support sns vpn m2 network-firewall storagegateway autoscaling servicecatalog fsx-lustre appsync snowmobile mediatailor ecr iottwinmaker elb codecommit memorydb lookoutvision s3 sms deepracer codeguruprofiler efs qldb mediaconvert fsx-openzfs comprehendmedical gamelift dynamodb docdb mediaconnect alexaforbusiness config cloudfront workmail securityhub appmesh deeplens cur sumerian rekognition lambda freertosota fsx-windows discovery ecs codebuild mediastore costexplorer]
```

### IsServiceInRegion

```
package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	res := awsservicemap.IsServiceInRegion("frauddetector", "eu-south-2")
	fmt.Println(res)
	res := awsservicemap.IsServiceInRegion("lambda", "us-east-1")
	fmt.Println(res)
}
```

Output: 
```
false
true
```

# Thanks

Thanks to Christophe Tafani-Dereeper (@christophetd) for the idea to make this a consumable library
