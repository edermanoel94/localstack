package localstack_test

import (
	"context"
	"github.com/edermanoel94/localstack"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {

	localStack, err := localstack.New(false)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, localStack)
}

func TestLocalStack_Run(t *testing.T) {

	t.Run("should dumping inspected container with make inspect enabled,", func(t *testing.T) {

		t.Cleanup(func() {
			os.Remove("dumping_inspected.json")
		})

		stack, err := localstack.New(true)

		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		err = stack.Run(ctx)

		assert.Nil(t, err)
		assert.FileExists(t, "dumping_inspected.json")
	})

	t.Run("should not dumping inspected container with make inspect disabled,", func(t *testing.T) {

		stack, err := localstack.New(false)

		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		err = stack.Run(ctx)

		assert.Nil(t, err)
		assert.NoFileExists(t, "dumping_inspected.json")
	})
}
