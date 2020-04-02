package localstack

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"strings"
)

type Service string

const (

	// SERVICE = "<name>/port", name of service in lowercase and port defined by localstack
	// the name of the services are listed in the aws tool: https://docs.aws.amazon.com/cli/latest/index.html
	// TODO: add more services

	S3    Service = "s3/4572"
	SNS           = "sns/4575"
	SQS           = "sqs/4576"
	Admin         = "admin/8080"
)

var all = []Service{S3, SNS, SQS, Admin}

func (s Service) Name() string {
	return strings.Split(string(s), "/")[0]
}

func (s Service) NatPort() nat.Port {
	return nat.Port(strings.Split(string(s), "/")[1])
}

func concatWithTCP(port nat.Port) nat.Port {
	return nat.Port(fmt.Sprintf("%s/tcp", port))
}
