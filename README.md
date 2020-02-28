# go-ssmhelpers

go-ssmhelpers is a library of useful functions to interact with the AWS Systems Manager (SSM) API via the AWS Golang SDK.
Used primarily for [ssm-helpers](github.com/disneystreaming/ssm-helpers) project.

https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/

## Build Requirements

*   go 1.12 or higher

## Golang Dependencies

Golang dependencies are managed by `go mod`.

### Test

```
go mod download
go test -v -cover ./...
```

## Install

```
go get "github.com/disneystreaming/go-ssmhelpers"
```
