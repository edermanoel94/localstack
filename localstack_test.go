package localstack_test

import (
	"context"
	"github.com/edermanoel94/localstack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	localStack, err := localstack.New()

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, localStack)
}

func TestLocalStack_Run(t *testing.T) {

	stack, err := localstack.New(localstack.S3, localstack.SNS)

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Run(ctx)

	assert.Nil(t, err)
}
