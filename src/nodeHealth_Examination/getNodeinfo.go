package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"time"
)
type StatusServer struct {
	Percent  StatusPercent
	CPU      []CPUInfo
	Mem      MemInfo
	Swap     SwapInfo
	Load     *load.AvgStat
	Network  map[string]InterfaceInfo
	BootTime uint64
	Uptime   uint64
}

type DockerStatu struct {
	DockerName   []string
	DockerStatus string
	DockerState	 string
}

type StatusPercent struct {
	CPU  float64
	Disk float64
	Mem  float64
	Swap float64
}
type CPUInfo struct {
	ModelName string
	Cores     int32
}
type MemInfo struct {
	Total     uint64
	Used      uint64
	Available uint64
}
type SwapInfo struct {
	Total     uint64
	Used      uint64
	Available uint64
}
type InterfaceInfo struct {
	Addrs    []string
	ByteSent uint64
	ByteRecv uint64
}

type DockerInfo struct {
	Name	    string
	Runtime		string
	status		string
}

func main() {
	//port := flag.String("port", ":8080", "HTTP listen port")
	//flag.Parse()
	//http.HandleFunc("/", getInfo)
	//err := http.ListenAndServe(*port, nil)
	//if err != nil {
	//	log.Fatalln("ListenAndServe: ", err)
	//}

	nI := ioutil.WriteFile("Node_Info.json", []byte(infoMiniJSON()), 0644)
	if nI != nil {
		fmt.Printf("WriteFileJson ERROR: %+v", nI)
	}

	dockerInfo()

	//dI := ioutil.WriteFile("Dokcer_Info.json", []byte(dockerInfo()), 0644)
	//if dI != nil {
	//	fmt.Printf("WriteFileJson ERROR: %+v", dI)
	//}


}
//func getInfo(w http.ResponseWriter, r *http.Request) {
////	w.Header().Set("Access-Control-Allow-Origin", "*")
////	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
////	w.Header().Set("Content-Type", "application/json")
////	fmt.Fprintf(w, infoMiniJSON())
////}

func infoMiniJSON() string {
	v, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	l, _ := load.Avg()
	i, _ := net.Interfaces()
	ss := new(StatusServer)
	ss.Load = l
	ss.Uptime = n.Uptime
	ss.BootTime = n.BootTime
	ss.Percent.Mem = v.UsedPercent
	ss.Percent.CPU = cc[0]
	ss.Percent.Swap = s.UsedPercent
	ss.Percent.Disk = d.UsedPercent
	ss.CPU = make([]CPUInfo, len(c))
	for i, ci := range c {
		ss.CPU[i].ModelName = ci.ModelName
		ss.CPU[i].Cores = ci.Cores
	}
	ss.Mem.Total = v.Total
	ss.Mem.Available = (v.Available /1024000)
	ss.Mem.Used = (v.Used /1024000)
	ss.Swap.Total = (s.Total /1024000)
	ss.Swap.Available = (s.Free /1024000)
	ss.Swap.Used = (s.Used /1024000)
	ss.Network = make(map[string]InterfaceInfo)
	for _, v := range nv {
		var ii InterfaceInfo
		ii.ByteSent = v.BytesSent
		ii.ByteRecv = v.BytesRecv
		ss.Network[v.Name] = ii
	}
	for _, v := range i {
		if ii, ok := ss.Network[v.Name]; ok {
			ii.Addrs = make([]string, len(v.Addrs))
			for i, vv := range v.Addrs {
				ii.Addrs[i] = vv.Addr
			}
			ss.Network[v.Name] = ii
		}
	}
	b, err := json.Marshal(ss)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}


func dockerInfo() string {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	//cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	di := new(DockerStatu)
	for _, container := range containers {
		di.DockerName = container.Names
		di.DockerState = container.State
		di.DockerStatus = container.Status
		js, err := json.Marshal(di)
		if err != nil {
			return ""
		} else {
			//dI := ioutil.WriteFile("Dokcer_Info.json", []byte(js), 0644)
			fd, _ := os.OpenFile("Dokcer_Info.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			json := []byte(js)
			fd.Write(json)
			//if dI != nil {
			//	fmt.Printf("WriteFileJson ERROR: %+v", dI)
			//}

		}
	}
	return ""
}