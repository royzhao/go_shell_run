package main

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"strconv"
	"strings"
)

type CpuSample struct {
	idle  uint64
	total uint64
}

func getMemUsage() float64 {
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct

	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	return float64(v.Total - v.Free)
}
func getCpuUsage(sample0, sample1 CpuSample) float64 {
	idleTicks := float64(sample1.idle - sample0.idle)
	totalTicks := float64(sample1.total - sample0.total)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
	return cpuUsage
}
func getCPUSample() (sample CpuSample) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				sample.total += val // tally up all the numbers to get total ticks
				if i == 4 {         // idle is the 5th field in the cpu line
					sample.idle = val
				}
			}
			return
		}
	}
	return
}

func getConsInfo() []Container {
	mutex.Lock()
	defer mutex.Unlock()
	num := len(mapCons)
	i := 0
	data := make([]Container, num)
	for _, c := range mapCons {
		data[i] = *c
		i++
	}
	return data
}
