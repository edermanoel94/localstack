# Localstack

A package which make easy initialize for localstack and debug your application with your best IDE, 
choosen services and without setup any Dockerfile to run. This working with network_mode=host

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
    
    // first parameter is to create a json file of localstack docker configurations and 
	// second parameter is service that you want to use, if not pass will get all services
    localStack, err := localstack.New(false, localstack.S3, localstack.SNS)

	if err != nil {
		// treat your error here
	}

	ctx := context.Background()

	err = localStack.Run(ctx)

	if err != nil {
		// treat your error here
	}
    
    // you can return the method Stop or Remove to put on defer and stop execution.
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
- [ ] Working with all network_modes
- [ ] Remove sleep and check services healthy
- [ ] Add config file for working with binary

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.
