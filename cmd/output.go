package cmd

import (
	"log"
	"fmt"
)


func Output(stdoutStderr []byte, err error) []byte {
	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
		log.Fatal(err)
	}
	return stdoutStderr
}
