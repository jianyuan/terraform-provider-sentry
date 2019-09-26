module github.com/jianyuan/terraform-provider-sentry

require (
	github.com/aws/aws-sdk-go v1.24.6 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform v0.12.4
	github.com/hashicorp/yamux v0.0.0-20180917205041-7221087c3d28 // indirect
	github.com/jianyuan/go-sentry v1.2.0
	github.com/mitchellh/go-homedir v1.1.0 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999 // indirect

go 1.13
