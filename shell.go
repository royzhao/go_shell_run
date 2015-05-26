package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

type IOHandler func(stdin io.Writer, stdout io.Reader, stderr io.Reader)

func System(cmd string, f IOHandler) (string, error) {
	c := exec.Command("/bin/bash", "-c", cmd)
	in, err := c.StdinPipe()
	if err != nil {
		return "error start cmd: " + cmd, err
	}
	e, err := c.StderrPipe()

	if err != nil {
		return "error start cmd: " + cmd, err
	}
	o, err := c.StdoutPipe()
	if err != nil {
		return "error start cmd: " + cmd, err
	}

	err = c.Start()
	log.Println("runing ", cmd)
	if err != nil {
		os.Exit(-1)
	}
	go f(in, o, e)

	err = c.Wait()
	//cancel timer
	log.Println("run command over")
	if err != nil {
		return "error", err
	}

	return "", nil
}
