package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var(
		cmd *exec.Cmd
		output []byte
		err error
	)
	//cmd = exec.Command("bash", "-c", "sleep 5;ls -l")
	cmd = exec.Command("/usr/bin/bash", "-c", "sleep 5;ls -l")
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}
