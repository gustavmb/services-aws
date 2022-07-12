/*
Simple UDS based client that connects to KVM socket file and set CPU Affinity
*/

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

type VMS struct {
	VMS []VM `json:"VM"`
}

type VM struct {
	SockPath string `json:sockpath`
	Vcpus    []int  `json:vcpus`
}

func reader(r io.Reader) string {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		log.Fatal(err)
	}
	return string(buf[0:n])
}

func connect(sock_path string) net.Conn {
	c, err := net.Dial("unix", sock_path)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func GetPIDs(serial_res string) []int {
	rx := regexp.MustCompile("CPU #[0-9]: thread_id=(.*?)\\r\\n")
	matches := rx.FindAllStringSubmatch(serial_res, -1)
	if matches == nil {
		log.Fatal("Bad output in console!")
	}
	var pids []int
	for _, match := range matches {
		intMatch, _ := strconv.Atoi(match[1]) //Get substring that contains PID only
		pids = append(pids, intMatch)
	}
	return pids
}

func SetCPUAffinity(cpu int, pid int) {
	var newMask unix.CPUSet
	newMask.Set(cpu)
	err := unix.SchedSetaffinity(pid, &newMask)
	if err != nil {
		log.Fatal("SchedSetaffinity: %v", err)
	}
	log.Println("PID: ", pid, "CPU: ", cpu, "Affinity Correctly Set")
}

func ReadFile() string {
	content, err := ioutil.ReadFile("sockdata/data")
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func main() {
	jsonFile, err := os.Open("sockdata/data")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	var vms VMS
	json.Unmarshal(byteValue, &vms)
	for _, vm := range vms.VMS {
		c := connect(vm.SockPath)
		defer c.Close()
		_, err := c.Write([]byte("info cpus\n"))
		if err != nil {
			log.Fatal("write error:", err)
		}
		time.Sleep(1e9)
		result := reader(c)
		pids := GetPIDs(result)
    if len(vm.Vcpus) != len(pids) {
      log.Fatal("Number of PIDS and CPUS are not equal for VM: ", vm.SockPath)
    }
		for idx, pid := range pids {
			go SetCPUAffinity(vm.Vcpus[idx], pid)
		}
	}
}
