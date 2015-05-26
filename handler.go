package main

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/nat"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Runres struct {
	cmd    string
	result string
}

func substring(str string, begin int, end int) (string, error) {
	r := []rune(str)
	if len(r) < end && begin < 0 {
		return "", newError(1, []byte("bad index"))
	}
	return string(r[begin:end]), nil
}
func getContainerPort(ports nat.PortMap) string {
	port_4470 := ports["4470/tcp"]
	if port_4470 != nil {
		return port_4470[0].HostPort
	}

	port_8080 := ports["8080/tcp"]
	if port_8080 != nil {
		return port_8080[0].HostPort
	}

	return ""
}

//docker inspect
func Handler4(stdin io.Writer, stdout io.Reader, stderr io.Reader) {
	obyte, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("err in exec handler4")
	} else {
		var res []types.ContainerJSON
		err = json.Unmarshal(obyte, &res)
		if err != nil {
			log.Println(err)
			return
		}
		// current := string(obyte)

		// res := strings.Fields(current)
		mutex.Lock()
		defer mutex.Unlock()
		cid, err := substring(res[0].Id, 0, 12)
		if err != nil {
			return
		}
		con := mapCons[cid]
		if con == nil {
			var con2 Container
			con2.Cpu = 0.0
			con2.Mem = 0.0
			con2.Id = cid
			con2.Port = getContainerPort(res[0].NetworkSettings.Ports)
			con2.Image = res[0].Config.Image
			mapCons[cid] = &con2
		} else {
			con.Image = res[0].Config.Image
		}
	}
}

//docker ps | grep {id}
func Handler3(stdin io.Writer, stdout io.Reader, stderr io.Reader) {
	obyte, err := ioutil.ReadAll(stdout)
	if err != nil {

	} else {
		current := string(obyte)
		res := strings.Fields(current)
		mutex.Lock()
		defer mutex.Unlock()

		con := mapCons[res[0]]
		if con == nil {
			var con2 Container
			con2.Cpu = 0.0
			con2.Mem = 0.0
			con2.Id = res[0]
			con2.Image = res[1]
			mapCons[res[0]] = &con2

		} else {
			con.Image = res[1]
		}
	}
}

//docker ps
func Handler2(stdin io.Writer, stdout io.Reader, stderr io.Reader) {
	obyte, err := ioutil.ReadAll(stdout)
	if err != nil {

	} else {
		current := string(obyte)
		tmp := strings.Split(current, "\n")
		num := len(tmp) - 1
		if num > 0 {
			mutex.Lock()
			defer mutex.Unlock()
			for i := 1; i < num; i++ {
				res := strings.Fields(tmp[i])
				log.Println(res)
				con := mapCons[res[0]]
				if con == nil {
					var con2 Container
					con2.Cpu = 0.0
					con2.Mem = 0.0
					con2.Id = res[0]
					con2.Image = res[1]
					mapCons[res[0]] = &con2

				} else {
					con.Image = res[1]
					con.Port = res[5]
				}
			}
		}
	}
}

//docker stats
func Handler1(stdin io.Writer, stdout io.Reader, stderr io.Reader) {
	for {

		obyte := make([]byte, 180)
		_, err := stdout.Read(obyte)
		if err != nil {
			log.Println(err)
		} else {
			current := string(obyte)
			tmp := strings.Split(current, "\n")
			if len(tmp) > 1 {
				res := strings.Fields(tmp[1])
				if len(res) > 6 {
					mutex.Lock()
					cpu := strings.Split(res[1], "%")

					cpu_data, err := strconv.ParseFloat(cpu[0], 64)
					mem_data, err := strconv.ParseFloat(res[2], 64)
					if err != nil {

					} else {
						con := mapCons[res[0]]
						if con == nil {
							var con2 Container
							con2.Cpu = cpu_data
							con2.Mem = 1024 * mem_data
							con2.Id = res[0]
							con = &con2
							mapCons[res[0]] = &con2
							System("docker inspect "+res[0], Handler4)
						} else {
							con.Cpu = cpu_data
							con.Mem = 1024 * mem_data
						}
						log.Println("==============")
						log.Println(*con)
						log.Println("==============")

					}
					mutex.Unlock()
				}
			}
		}

	}
}
