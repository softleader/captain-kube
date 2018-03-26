package log

import (
	"fmt"
	"log"
)

func Output(stdoutStderr []byte, err error) []byte {
	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
		log.Fatal(err)
	}
	return stdoutStderr
}

func Command(command string) {
	fmt.Println(">", command)
}
