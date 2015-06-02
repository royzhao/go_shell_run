package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	data      SendData
	cons      *[]Container
	mapCons   = make(map[string]*Container)
	updateMap = make(map[string]int64)
	lastCpu   CpuSample
	mutex     sync.Mutex
)

func main() {
	if len(os.Args) < 2 {
		log.Println("enter the peer")
		os.Exit(-1)
	}
	peer := os.Args[1]
	//get first cpu sample
	lastCpu = getCPUSample()
	//inint client
	client, err := NewPeerClient(peer)
	if err != nil {
		os.Exit(-1)
	}

	// run container
	// go func() {
	// 	r := Runres{
	// 		cmd:    "docker stats $(docker ps -q)",
	// 		result: "",
	// 	}
	// 	_, err := System(r.cmd, Handler1)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		_, err = System("docker ps", Handler2)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }()

	for {
		time.Sleep(1 * time.Second)
		var sd SendData
		_, err = System("docker ps", Handler2)
		if err != nil {
			log.Println(err)
		} else {
			newCpu := getCPUSample()
			sd.Cpu = getCpuUsage(lastCpu, newCpu)
			sd.Mem = getMemUsage()
			sd.Containers = getConsInfo()
			log.Println("run..")

			// log.Println(sd)
			client.Send(sd)
		}
	}

}
