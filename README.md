# aws-servicemap
Not all AWS services are supported in all regions. Use this is a go module to return either a slice of supported regions for a service, or a slice of supported services for a region.   


## Example Usage

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