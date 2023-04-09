package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	// go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	stdinPip, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(stdinPip, "23 11\n")
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(out.String())
}
