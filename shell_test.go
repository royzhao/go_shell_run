package main

import (
	"testing"
)

// func Test_dockerps(t *testing.T) {
// 	_, err := System("docker inspect 2a9be08a603b", Handler4)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func Test_dockertop(t *testing.T) {
	_, err := System("docker stats $(docker ps -q)", Handler1)
	if err != nil {
		t.Error(err)
	}
}
