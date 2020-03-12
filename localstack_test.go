package localstack_test

import (
	"context"
	"github.com/edermanoel94/localstack-api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLocalStack(t *testing.T) {

	localStack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, localStack)
}

func TestLocalStack_Pull(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Pull(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStack_Create(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Create(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStack_Start(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Start(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStack_Stop(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Stop(ctx, nil)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStack_Remove(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Remove(ctx, true)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStack_LocalStackContainerExist(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	exists, err := stack.ContainerExists(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, exists)
}

func TestLocalStack_Run(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = stack.Run(ctx)

	if err != nil {
		t.Fatal(err)
	}

}

func TestLocalStack_IsRunning(t *testing.T) {

	stack, err := localstack.NewLocalStack()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	isRunning, err := stack.IsRunning(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, false, isRunning)
}

func TestLocalStack_Logs(t *testing.T) {

}
