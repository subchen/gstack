package main

import (
	"fmt"
	"time"

	"github.com/subchen/gstack/errors"
)

func createFile() error {
	return errors.New("file not permission")
}

func writeFile() error {
	err := createFile()
	if err != nil {
		return errors.Wrap(err, "file write error")
	}

	return nil
}

func main() {
	err := writeFile()
	fmt.Printf("error %%v: %v\n\n", err)
	fmt.Printf("error %%+v: %+v\n\n", err)

	cause := errors.Cause(err)
	fmt.Printf("cause %%v: %v\n\n", cause)
	fmt.Printf("cause %%+v: %+v\n\n", cause)

	go func() {
		err := writeFile()
		fmt.Printf("go error %%v: %v\n\n", err)
		fmt.Printf("go error %%+v: %+v\n\n", err)
	}()

	time.Sleep(1 * time.Second)
}
