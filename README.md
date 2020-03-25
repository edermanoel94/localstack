# Localstack

A package which make easy initialize for localstack and debug your application with your best IDE, 
choosen services and without setup any Dockerfile to run.

## Example

```go
package your_package_test

import (
    "context"
    "github.com/edermanoel94/localstack"
    "log"
    "testing"
)

func before() {
    
    // pass services which u will use
    localStack, err := localstack.New(localstack.S3, localstack.SQS)

	if err != nil {
		// treat your error here
	}

	ctx := context.Background()

	err = localStack.Run(ctx)

	if err != nil {
		// treat your error here
	}
}

func TestSomeHandler(t *testing.T) {
    
    before()
    
    // reference your aws-sdk to localstack

    // do your stuff with your test here
}
```

## Prerequisites

- Go 1.14
- Docker

## Installation

We using go modules and get always the last version

```
$ go get github.com/edermanoel94/localstack@latest
```

# TODO List

- [ ] Add more services
- [ ] Better costumization on container
- [ ] Create helper for aws-sdk

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.