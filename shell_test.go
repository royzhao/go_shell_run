package main

import (
	"log"
	"testing"
	"time"
)

// func Test_dockerps(t *testing.T) {
// 	_, err := System("docker inspect 2a9be08a603b", Handler4)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func Test_dockertop(t *testing.T) {
	// go func() {
	for {
		time.Sleep(4 * time.Second)
		ddd := getConsInfo()
		for _, con2 := range ddd {
			// log.Println("id:", con2)
			// if con2 != nil {
			log.Println("name:", con2.Image, " cpu:", con2.Cpu, " id", con2.Id, " mem:", con2.Mem, " port:", con2.Port)

			// }
		}
		_, err := System("docker ps", Handler2)
		if err != nil {
			t.Error(err)
		}
	}
	// }()
	// _, err := System("docker stats $(docker ps -q)", Handler1)
	// if err != nil {
	// 	t.Error(err)
	// }
}
