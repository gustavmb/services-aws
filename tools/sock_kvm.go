/*
Simple UDS based client that connects to KVM socket file and set CPU Affinity
*/

package main

import (
        "encoding/json"
        "io"
        "os"
        "io/ioutil"
        "log"
        "net"
        "time"
        "regexp"
        "fmt"
        "strconv"
        "golang.org/x/sys/unix"
)

type VMS struct {
    VMS []VM `json:"VM"`
}

type VM struct {
  SockPath string `json:sockpath`
  Vcpus []int `json:vcpus`

}


func reader (r io.Reader) (string, error) {
        buf := make([]byte, 1024)
        n, err := r.Read(buf[:])
        if err != nil {
                return "", err
        }
        return string(buf[0:n]), err
}

func connect (sock_path string) net.Conn {
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
    for _, match:= range matches {
      intMatch, _ := strconv.Atoi(match[1]) //Get substring that contains PID only
      pids = append(pids, intMatch)
    }
    return pids
}

func SetCPUAffinity(cpus []int, pids []int) {
    if len(cpus) != len(pids) {
        log.Fatal("Number of PIDS and CPUS does not coincide")
    }
    for idx, pid:= range pids {
        var newMask unix.CPUSet
        newMask.Set(cpus[idx])
        err := unix.SchedSetaffinity(pid, &newMask)
        if err != nil {
            fmt.Printf("SchedSetaffinity: %v", err)
        }

    }
}

func ReadFile() string {
    content,err := ioutil.ReadFile("sockdata/data")
    if err != nil {
        log.Fatal(err)
    }
    return string(content)
}

func main() {
        jsonFile, err := os.Open("sockdata/data")
        if err != nil {
            fmt.Println(err)
        }
        defer jsonFile.Close()
        byteValue, _ := ioutil.ReadAll(jsonFile)
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
            result, err := reader(c)
            pids := GetPIDs(result)
            SetCPUAffinity(vm.Vcpus, pids)
        }
}
