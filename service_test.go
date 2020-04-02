package localstack_test

import (
	"github.com/edermanoel94/localstack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Name(t *testing.T) {

	testCases := []struct {
		description string
		service     localstack.Service
		expect      string
	}{
		{"should get name of service S3", localstack.S3, "s3"},
		{"should get name of service SNS", localstack.SNS, "sns"},
		{"should get name of service SQS", localstack.SQS, "sqs"},
	}

	for _, tc := range testCases {

		t.Run(tc.description, func(t *testing.T) {

			assert.Equal(t, tc.expect, tc.service.Name())
		})
	}
}

func TestService_NatPort(t *testing.T) {

	testCases := []struct {
		description string
		service     localstack.Service
		expect      string
	}{
		{"should get port of service S3", localstack.S3, "4572"},
		{"should get port of service SNS", localstack.SNS, "4575"},
		{"should get port of service SQS", localstack.SQS, "4576"},
	}

	for _, tc := range testCases {

		t.Run(tc.description, func(t *testing.T) {

			assert.Equal(t, tc.expect, tc.service.NatPort().Port())
		})
	}
}
